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
			_ = formatter.RespJSON(500, map[string]string{"error": "Internal Server Error"}, w)
			return
		}
		resp, err := client.Register(ctx, req.Email, req.Name, req.Password)
		if err != nil {
			_ = formatter.RespJSON(500, map[string]string{"error": "Internal Server Error"}, w)
			return
		}

		err = formatter.RespJSON(200, dto.ResponseRegisterUser{UserId: resp.UserId}, w)
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
			_ = formatter.RespJSON(400, map[string]string{"error": "Bad Request"}, w)
			return
		}

		resp, err := authClient.Login(ctx, req.Email, req.Password)

		if err != nil {
			_ = formatter.RespJSON(500, map[string]string{"error": "Internal Server Error"}, w)
			return
		}

		err = formatter.RespJSON(200, dto.ResponseLoginUser{
			AccessToken:  resp.AccessToken,
			RefreshToken: "", //TODO: add refresh token
			UserId:       resp.UserId,
		}, w)
		if err != nil {
			log.Println(err)
		}

	}
}
