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
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			accessToken := parts[1]

			isValid, _ := authClient.Validate(ctx, accessToken)

			if !isValid {
				refreshCookie, err := r.Cookie("refresh_token")
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				resp, err := authClient.Refresh(ctx, refreshCookie.Value)
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "access_token",
					Value:    resp.AccessToken,
					HttpOnly: true,
					Path:     "/",
				})
			}

			next.ServeHTTP(w, r)
		})
	}
}
