package auth

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

var ErrLDAPUserNotFound = errors.New("ldap user not found")

var reDigitsOnly = regexp.MustCompile(`^\d+$`)

// escapeDN escapes a string for use in an LDAP DN per RFC 4514.
func escapeDN(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == 0:
			b.WriteString(`\00`)
		case i == 0 && (c == ' ' || c == '#'):
			fmt.Fprintf(&b, `\%02x`, c)
		case i == len(s)-1 && c == ' ':
			fmt.Fprintf(&b, `\%02x`, c)
		case c == ',' || c == '+' || c == '"' || c == '\\' ||
			c == '<' || c == '>' || c == ';' || c == '=':
			b.WriteByte('\\')
			b.WriteByte(c)
		case c < 0x20:
			fmt.Fprintf(&b, `\%02x`, c)
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

type LDAPAuth struct {
	url           string
	baseDN        string
	serviceDN     string
	servicePasswd string
	useStartTLS   bool
	tlsConfig     *tls.Config
}

func NewLDAPAuth(url, baseDN, servicePasswd string) (*LDAPAuth, error) {
	tlsCfg := &tls.Config{InsecureSkipVerify: false, MinVersion: tls.VersionTLS12}

	if strings.HasPrefix(url, "ldaps://") {
		return &LDAPAuth{
			url:           url,
			baseDN:        baseDN,
			serviceDN:     fmt.Sprintf("cn=svc-localis,%s", baseDN),
			servicePasswd: servicePasswd,
			useStartTLS:   false,
			tlsConfig:     tlsCfg,
		}, nil
	}

	conn, err := ldap.DialURL(url, ldap.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}))
	if err != nil {
		return nil, fmt.Errorf("ldap dial: %w", err)
	}
	defer conn.Close()

	useStartTLS := false
	if strings.HasPrefix(url, "ldap://") {
		if err := conn.StartTLS(tlsCfg); err != nil {
			slog.Warn("LDAP StartTLS failed, continuing without TLS", "error", err)
		} else {
			useStartTLS = true
		}
	}

	return &LDAPAuth{
		url:           url,
		baseDN:        baseDN,
		serviceDN:     fmt.Sprintf("cn=svc-localis,%s", baseDN),
		servicePasswd: servicePasswd,
		useStartTLS:   useStartTLS,
		tlsConfig:     tlsCfg,
	}, nil
}

type LDAPUser struct {
	EmployeeNumber string
	Name           string
	Email          string
	Groups         []string
}

func (a *LDAPAuth) Bind(employeeNumber, password string) (*LDAPUser, error) {
	conn, err := ldap.DialURL(a.url, ldap.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}))
	if err != nil {
		return nil, fmt.Errorf("ldap dial: %w", err)
	}
	defer conn.Close()

	if a.useStartTLS {
		if err := conn.StartTLS(a.tlsConfig); err != nil {
			return nil, fmt.Errorf("ldap starttls: %w", err)
		}
	}

	err = conn.Bind(a.serviceDN, a.servicePasswd)
	if err != nil {
		return nil, fmt.Errorf("ldap service bind: %w", err)
	}

	if !reDigitsOnly.MatchString(employeeNumber) {
		return nil, ErrLDAPUserNotFound
	}

	searchReq := ldap.NewSearchRequest(
		a.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(objectClass=posixAccount)(uidNumber=%s))", ldap.EscapeFilter(employeeNumber)),
		[]string{"uid", "uidNumber", "cn", "mail", "memberOf", "gidNumber", "primaryGroup"},
		nil,
	)

	result, err := conn.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("ldap search: %w", err)
	}

	if len(result.Entries) == 0 {
		return nil, ErrLDAPUserNotFound
	}

	entry := result.Entries[0]
	userCN := entry.GetAttributeValue("cn")
	escapedCN := escapeDN(userCN)

	// Search for groups while still authenticated as service account
	gidNumber := entry.GetAttributeValue("gidNumber")
	if gidNumber == "" {
		gidNumber = entry.GetAttributeValue("primaryGroup")
	}

	var groups []string
	if gidNumber != "" {
		groupSearchReq := ldap.NewSearchRequest(
			a.baseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0, 0, false,
			fmt.Sprintf("(&(objectClass=posixGroup)(gidNumber=%s))", ldap.EscapeFilter(gidNumber)),
			[]string{"cn"},
			nil,
		)

		groupResult, err := conn.Search(groupSearchReq)
		if err != nil {
			slog.Warn("ldap group search failed, defaulting to empty groups", "gid_number", gidNumber, "error", err)
		} else {
			for _, g := range groupResult.Entries {
				groups = append(groups, g.GetAttributeValue("cn"))
			}
		}
	}

	// Now bind as user to verify password
	userDN := fmt.Sprintf("cn=%s,%s", escapedCN, a.baseDN)
	err = conn.Bind(userDN, password)
	if err != nil {
		return nil, fmt.Errorf("ldap bind: %w", err)
	}

	user := &LDAPUser{
		EmployeeNumber: entry.GetAttributeValue("uidNumber"),
		Name:           userCN,
		Email:          entry.GetAttributeValue("mail"),
		Groups:         groups,
	}

	return user, nil
}

func (a *LDAPAuth) ResolveRole(user *LDAPUser) string {
	for _, g := range user.Groups {
		switch strings.ToLower(g) {
		case "admin", "administrators", "admins":
			return "admin"
		case "manager", "managers":
			return "manager"
		case "supervisor", "supervisors":
			return "supervisor"
		}
	}
	return "technician"
}
