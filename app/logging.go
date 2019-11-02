package app

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

var Logging = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"source":    r.RemoteAddr,
			"uri":       r.RequestURI,
			"useragent": r.UserAgent(),
			"method":    r.Method,
		}).Trace("Handling Request")
		next.ServeHTTP(w, r)
	})
}
