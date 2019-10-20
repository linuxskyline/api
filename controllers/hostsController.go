package controllers

import (
	"apt-api/models"
	u "apt-api/utils"
	"encoding/json"
	"net/http"
)

var CreateHost = func(w http.ResponseWriter, r *http.Request) {
	host := &models.Host{}

	err := json.NewDecoder(r.Body).Decode(host)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	resp := host.Create()
	u.Respond(w, resp)
}

var ListHosts = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetHosts()
	resp := u.Message(true, "success")

	for _, host := range data {
		host.SecurityUpdates = host.CountSecurityUpdates()
		host.Updates = host.CountUpdates()
	}

	resp["data"] = data
	u.Respond(w, resp)
}
