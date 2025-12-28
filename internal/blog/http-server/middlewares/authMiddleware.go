package middlewares

import (
	"habr/internal/blog/grpc/client"
	"net/http"
	"strings"
)

func AuthMiddleware(authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				SendUnauthorized(w)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				SendUnauthorized(w)
			}

			accessToken := parts[1]

			// возвращает bool и int32
			// error проверяется внутри функции
			isValid, err := authClient.Validate(ctx, accessToken)
			if err != nil {
				SendUnauthorized(w)
			}
			if !isValid {
				refreshCookie, err := r.Cookie("refresh_token")
				if err != nil {
					SendUnauthorized(w)
				}

				resp, err := authClient.Refresh(ctx, refreshCookie.Value)
				if err != nil {
					SendUnauthorized(w)
				}

				// Отправляем новый access token клиенту через заголовок ответа
				w.Header().Set("X-New-Access-Token", resp.AccessToken)
				SendUnauthorized(w)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func SendUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	return
}
