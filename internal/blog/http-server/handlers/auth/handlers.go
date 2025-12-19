package auth

import (
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
		resp, err := client.Register(ctx, req.Email, req.Name, req.Password)
		if err != nil {
			_ = formatter.RespJSON(500, map[string]string{"error": "Internal Server Error"}, w)
			return
		}

		err = formatter.RespJSON(200, dto.ResponseRegisterUser{Id: resp.UserId}, w)
		if err != nil {
			log.Println(err)
		}
	}
}

func LoginUser(authClient *client.AuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
