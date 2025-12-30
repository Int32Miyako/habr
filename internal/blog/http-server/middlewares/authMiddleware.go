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
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				SendUnauthorized(w)
				return
			}

			accessToken := parts[1]

			// возвращает bool и int32
			// error проверяется внутри функции
			isValid, err := authClient.Validate(ctx, accessToken)
			if err != nil {
				SendUnauthorized(w)
				return
			}
			if !isValid {
				refreshCookie, err := r.Cookie("refresh_token")
				if err != nil {
					SendUnauthorized(w)
					return
				}

				resp, err := authClient.Refresh(ctx, refreshCookie.Value)
				if err != nil {
					SendUnauthorized(w)
					return
				}

				// Отправляем новый access token клиенту через заголовок ответа
				w.Header().Set("X-New-Access-Token", resp.AccessToken)
				SendUnauthorized(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func SendUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
