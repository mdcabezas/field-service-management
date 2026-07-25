package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type JWTClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
}

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		segments := strings.Split(tokenStr, ".")
		if len(segments) != 3 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Verify HMAC-SHA256 signature
		signingInput := segments[0] + "." + segments[1]
		expectedSig := hmacSHA256([]byte(signingInput), []byte(secret))
		actualSig, err := base64URLDecode(segments[2])
		if err != nil {
			http.Error(w, "Invalid token signature encoding", http.StatusUnauthorized)
			return
		}

		if !hmac.Equal(expectedSig, actualSig) {
			http.Error(w, "Invalid token signature", http.StatusUnauthorized)
			return
		}

		// Decode payload
		payload, err := base64URLDecode(segments[1])
		if err != nil {
			http.Error(w, "Invalid token payload encoding", http.StatusUnauthorized)
			return
		}

		var claims JWTClaims
		if err := json.Unmarshal(payload, &claims); err != nil {
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}

		// Check expiration
		if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// Return headers for Traefik to inject into downstream request
		w.Header().Set("X-User-ID", claims.Sub)
		w.Header().Set("X-User-Role", claims.Role)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("JWT Validator listening on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func hmacSHA256(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func base64URLDecode(s string) ([]byte, error) {
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}
