package middlewares

import (
	"habr/internal/blog/grpc/client"
	"net/http"
)

func AuthMiddleware(authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пример проверки токена
			ctx := r.Context()
			token := r.Header.Get("Authorization")
			if token == "" || !authClient.Validate(ctx, token) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
