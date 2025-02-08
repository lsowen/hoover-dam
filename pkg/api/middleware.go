package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lsowen/hoover-dam/pkg/config"
)

func AuthMiddleware(cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("authorization")
			if authHeader == "" {
				http.Error(w, "No `Authorization` header", http.StatusForbidden)
				return
			}

			segments := strings.Fields(authHeader)
			if len(segments) != 2 || !strings.EqualFold(segments[0], "bearer") {
				http.Error(w, "Malformed `Authorization` header", http.StatusForbidden)
				return
			}

			_, err := jwt.Parse(segments[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.Auth.Encrypt.SecretKey), nil
			}, jwt.WithAudience("auth-client"), jwt.WithExpirationRequired())
			if err != nil {
				http.Error(w, "Failed to validate bearer token", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
