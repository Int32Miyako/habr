package middlewares

import (
	"context"
	"habr/internal/blog/grpc/client"
	"net/http"
	"strings"
)

func AuthMiddleware(authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 🔴 1. Достаём заголовок
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// 🔴 2. Проверяем формат
			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, prefix)

			// 🔴 3. Валидируем токен
			claims, err := validateToken(r.Context(), token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// 🟢 4. Кладём claims в context
			ctx := context.WithValue(r.Context(), "user", claims)

			// 🟢 5. Пропускаем дальше
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateToken(ctx context.Context, token string) (interface{}, interface{}) {
	return ctx, token
}
