package auth

import (
	"encoding/json"
	"habr/internal/auth/app/http/dto"
	"habr/internal/auth/app/services"
	"habr/internal/pkg/formatter"
	"log"
	"net/http"
)

func RegisterUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := dto.RequestRegisterUser{}

		if req.Email == "" {
			_ = formatter.RespBadRequest("Bad Request", w)
			return
		}
		if req.Username == "" {
			_ = formatter.RespBadRequest("Bad Request", w)
			return
		}
		if req.Password == "" {
			_ = formatter.RespBadRequest("Bad Request", w)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			_ = formatter.RespBadRequest("Bad Request", w)
			return
		}

		userId, err := userService.RegisterUser(ctx, req.Email, req.Username, req.Password)
		if err != nil {
			status := formatter.ErrorToStatus(err)
			_ = formatter.RespError(status, err.Error(), w)

			return
		}

		err = formatter.RespJSON(http.StatusCreated, dto.ResponseRegisterUser{UserId: userId}, w)
		if err != nil {
			log.Println(err)
		}
	}
}
