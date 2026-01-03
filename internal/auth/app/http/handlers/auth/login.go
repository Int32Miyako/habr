package auth

import (
	"encoding/json"
	"habr/internal/auth/app/services"
	"habr/internal/blog/http/dto"
	"habr/internal/pkg/formatter"
	"log"
	"net/http"
)

func LoginUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := dto.RequestLoginUser{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			_ = formatter.RespBadRequest("Bad Request", w)
			return
		}

		resp, err := userService.LoginUser(ctx, dto.RequestLoginUser{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			status := formatter.ErrorToStatus(err)
			_ = formatter.RespError(status, err.Error(), w)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			Secure:   true, // обязательно в проде (https)
		})

		err = formatter.RespJSON(http.StatusOK, dto.ResponseLoginUser{
			AccessToken: resp.AccessToken,
			UserId:      resp.UserId,
		}, w)
		if err != nil {
			log.Println(err)
		}
	}
}
