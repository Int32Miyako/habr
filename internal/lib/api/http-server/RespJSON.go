package http_server

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespJSON(w http.ResponseWriter, obj any) {
	w.Header().Set("Content-Type", "application/json")
	m, err := json.Marshal(obj)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = w.Write(m)
	if err != nil {
		log.Fatal(err)
		return
	}
}
