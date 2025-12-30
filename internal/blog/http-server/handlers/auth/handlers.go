package auth

import (
	"encoding/json"
	"habr/internal/blog/grpc/client"
	"habr/internal/blog/http-server/dto"
	"habr/internal/pkg/formatter"
	"log"
	"net/http"
)

func RegisterUser(client *client.AuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := dto.RequestRegisterUser{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			_ = formatter.RespBadRequest("Bad Request", w)
			return
		}
		resp, err := client.Register(ctx, req.Email, req.Name, req.Password)
		if err != nil {
			status := formatter.ErrorToStatus(err)
			_ = formatter.RespError(status, err.Error(), w)
			return
		}

		err = formatter.RespJSON(http.StatusCreated, dto.ResponseRegisterUser{UserId: resp.UserId}, w)
		if err != nil {
			log.Println(err)
		}
	}
}

func LoginUser(authClient *client.AuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := dto.RequestLoginUser{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			_ = formatter.RespBadRequest("Bad Request", w)
			return
		}

		resp, err := authClient.Login(ctx, req.Email, req.Password)

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
