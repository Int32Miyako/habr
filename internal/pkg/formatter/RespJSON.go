package formatter

import (
	"encoding/json"
	"errors"
	"habr/internal/auth/core/constants"
	"log"
	"net/http"
)

func RespJSON(status int, obj any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, `{"error":"failed to marshal response"}`, http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func RespError(status int, message string, w http.ResponseWriter) error {
	return RespJSON(status, map[string]string{"error": message}, w)
}

func RespInternalError(w http.ResponseWriter) error {
	return RespError(http.StatusInternalServerError, "Internal Server Error", w)
}

func RespBadRequest(message string, w http.ResponseWriter) error {
	return RespError(http.StatusBadRequest, message, w)
}

func RespOK[T any](data T, w http.ResponseWriter) error {
	return RespJSON(http.StatusOK, data, w)
}

func RespUnauthorized(message string, w http.ResponseWriter) error {
	return RespError(http.StatusUnauthorized, message, w)
}

func ErrorToStatus(err error) int {
	switch {
	case errors.Is(err, constants.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, constants.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, constants.ErrUserAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
