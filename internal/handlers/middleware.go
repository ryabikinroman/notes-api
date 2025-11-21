package handlers

import (
	"context"
	"net/http"
	"notes-api/internal/auth"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			WriteError(w, http.StatusUnauthorized, "Токен отсутствует")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			WriteError(w, http.StatusUnauthorized, "Неверный формат токена")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &auth.Claims{}

		secret := []byte(os.Getenv("JWT_SECRET"))
		if len(secret) == 0 {
			WriteError(w, http.StatusInternalServerError, "JWT секрет отсутствует")
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			WriteError(w, http.StatusUnauthorized, "Неверный или просроченный токен")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
