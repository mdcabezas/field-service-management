package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrJWTSecretEmpty    = errors.New("jwt secret cannot be empty")
	ErrJWTSecretTooShort = errors.New("jwt secret must be at least 32 bytes")
	ErrInvalidToken      = errors.New("invalid token")
	ErrWrongTokenType    = errors.New("wrong token type")
)

type JWTAuth struct {
	secret []byte
}

func NewJWTAuth(secret string) (*JWTAuth, error) {
	if secret == "" {
		return nil, ErrJWTSecretEmpty
	}
	if len(secret) < 32 {
		return nil, fmt.Errorf("%w: got %d bytes", ErrJWTSecretTooShort, len(secret))
	}

	return &JWTAuth{
		secret: []byte(secret),
	}, nil
}

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	UserID    string `json:"sub"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func (a *JWTAuth) GenerateToken(userID, role, tokenType string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "localis-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

func (a *JWTAuth) ValidateToken(tokenString string, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secret, nil
	}, jwt.WithIssuer("localis-backend"))

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if expectedType != "" && claims.TokenType != expectedType {
		return nil, ErrWrongTokenType
	}

	return claims, nil
}
