package formatter

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespJSON(status int, obj any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	data, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, `{"error":"failed to marshal response"}`, http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
