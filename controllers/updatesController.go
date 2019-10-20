package controllers

import (
	"apt-api/models"
	u "apt-api/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var CreateUpdate = func(w http.ResponseWriter, r *http.Request) {
	update := &models.Update{}
	var hostID int

	if r.Context().Value("authed") != true {
		hostToken := r.Header.Get("HostToken")

		host, err := models.GetHostMatchingToken(hostToken)
		if err != nil {

			u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
			return
		}

		hostID = int(host.ID)
	} else {
		var err error
		hostID, err = strconv.Atoi(r.URL.Query()["id"][0])
		if err != nil {
			u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
			return
		}
	}

	err := json.NewDecoder(r.Body).Decode(update)
	if err != nil {
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request: %s", err)))
		return
	}

	update.HostId = uint(hostID)

	resp := update.Create()
	u.Respond(w, resp)
}

var ListUpdates = func(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		//The passed path parameter is not an integer
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request: %s", err)))
		return
	}

	data := models.GetUpdates(uint(id))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
