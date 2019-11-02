package controllers

import (
	"apt-api/models"
	u "apt-api/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/google/jsonapi"
)

func authAndGetId(r *http.Request) (int, error) {
	var hostID int

	if r.Context().Value("authed") != true {
		hostToken := r.Header.Get("HostToken")

		host, err := models.GetHostMatchingToken(hostToken)
		if err != nil {
			return 0, err
		}

		hostID = int(host.ID)
	} else {
		var err error
		hostID, err = strconv.Atoi(r.URL.Query()["id"][0])
		if err != nil {
			return 0, err
		}
	}

	return hostID, nil
}

var CreateUpdate = func(w http.ResponseWriter, r *http.Request) {
	update := &models.Update{}

	hostID, err := authAndGetId(r)
	if err != nil {
		log.WithFields(log.Fields{
			"source":    r.RemoteAddr,
			"uri":       r.RequestURI,
			"useragent": r.UserAgent(),
			"cause":     err,
		}).Trace("Failed to auth")
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
	}

	err = json.NewDecoder(r.Body).Decode(update)
	if err != nil {
		log.WithFields(log.Fields{
			"source":    r.RemoteAddr,
			"uri":       r.RequestURI,
			"useragent": r.UserAgent(),
			"cause":     err,
		}).Trace("Failed to decode body")
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request: %s", err)))
		return
	}

	update.HostId = uint(hostID)

	resp := update.Create()
	u.Respond(w, resp)
}

var ListUpdates = func(w http.ResponseWriter, r *http.Request) {
	hostID, err := authAndGetId(r)
	if err != nil {
		log.WithFields(log.Fields{
			"source":    r.RemoteAddr,
			"uri":       r.RequestURI,
			"useragent": r.UserAgent(),
			"cause":     err,
		}).Trace("Failed to auth")
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
	}

	data := models.GetUpdates(uint(hostID))

	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := jsonapi.MarshalPayload(w, data); err != nil {
		log.Error("Failed marshalling payload")
	}
}

var DeleteUpdate = func(w http.ResponseWriter, r *http.Request) {
	hostID, err := authAndGetId(r)
	if err != nil {
		log.WithFields(log.Fields{
			"source":    r.RemoteAddr,
			"uri":       r.RequestURI,
			"useragent": r.UserAgent(),
			"cause":     err,
		}).Warning("Failed to auth")
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
	}

	vars := mux.Vars(r)


	updateID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"source":    r.RemoteAddr,
			"uri":       r.RequestURI,
			"useragent": r.UserAgent(),
			"providedid": vars["id"],
			"cause":     err,
		}).Error("Failed decode update id from url")
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
	}

	log.WithFields(log.Fields{
		"id": updateID,
	}).Trace("Deleting host")

	data := models.GetUpdate(uint(updateID))

	if data.HostId != uint(hostID) {
		log.WithFields(log.Fields{
			"source":     r.RemoteAddr,
			"uri":        r.RequestURI,
			"useragent":  r.UserAgent(),
			"authedhost": hostID,
			"updatehost": data.HostId,
		}).Error("Attempted deletion of update for wrong host")
		u.Respond(w, u.Message(false, fmt.Sprintf("There was an error in your request")))
	}

	u.Respond(w, data.Delete())
}
