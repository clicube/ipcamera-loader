package main

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
}

func (h *Handler) errorMessage(message string) []byte {
	body := map[string]interface{}{
		"result":  "NG",
		"message": message,
	}
	bytes, _ := json.Marshal(body)
	return bytes
}

func (h *Handler) responseJson(writer http.ResponseWriter, body interface{}) {
	bytes, _ := json.Marshal(body)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(bytes)
}
