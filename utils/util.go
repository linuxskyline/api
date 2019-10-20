package utils

import (
	"encoding/json"
	"net/http"
)

type ReturnMessage map[string]interface{}

func Message(status bool, message string) ReturnMessage {
	return ReturnMessage{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data ReturnMessage) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(data)
}
