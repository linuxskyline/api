package controllers

import (
	"apt-api/models"
	u "apt-api/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

var GetHost = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp := u.Message(true, "success")

	hostID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"source":     r.RemoteAddr,
			"uri":        r.RequestURI,
			"useragent":  r.UserAgent(),
			"providedid": vars["id"],
			"cause":      err,
		}).Error("Failed decode update id from url")
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
	}

	host := models.GetHost(uint(hostID))

	resp["data"] = host
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

var DeleteHost = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hostID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"source":     r.RemoteAddr,
			"uri":        r.RequestURI,
			"useragent":  r.UserAgent(),
			"providedid": vars["id"],
			"cause":      err,
		}).Error("Failed decode update id from url")
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
	}

	log.WithFields(log.Fields{
		"id": hostID,
	}).Trace("Deleting host")

	data := models.GetHost(uint(hostID))

	u.Respond(w, data.Delete())
}
