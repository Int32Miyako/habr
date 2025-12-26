package middlewares

import (
	"habr/internal/blog/grpc/client"
	"net/http"
)

func AuthMiddleware(authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accessToken := r.Header.Get("Authorization")
			if accessToken == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if !authClient.Validate(ctx, accessToken) {
				cookie, err := r.Cookie("refresh_token")
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				resp, err := authClient.Refresh(ctx, cookie.Value)
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				w.Header().Set("X-New-Access-Token", resp.AccessToken)
				r.Header.Set("Authorization", resp.AccessToken)
			}
			next.ServeHTTP(w, r)
		})
	}
}
