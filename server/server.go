package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/makpoc/hades-sheet/gsclient"
)

var sheet *gsclient.Sheet

func initSheet() (*gsclient.Sheet, error) {
	spreadsheetId, ok := os.LookupEnv("SHEET_ID")
	if !ok {
		return nil, fmt.Errorf("No SHEET_ID found in environment")
	}
	var err error
	sheet, err = gsclient.New(spreadsheetId)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet client: %v", err)
	}

	return sheet, nil
}

func Start() error {
	router := mux.NewRouter().StrictSlash(true)
	s := router.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/timezones", auth(timeLogger(timeZonesHandler)))
	s.HandleFunc("/users", auth(timeLogger(usersHandler)))

	var err error
	sheet, err = initSheet()
	if err != nil {
		return err
	}

	log.Println("Starting server...")
	return http.ListenAndServe(":3000", s)
}

// auth provides uthorization layer based on secret in query parameter
func auth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var secret string
		var ok bool
		if secret, ok = os.LookupEnv("secret"); !ok {
			// server was not configured with secret
			h.ServeHTTP(w, r)
			return
		}

		querySecret := r.URL.Query().Get("secret")
		if querySecret != "" && querySecret == secret {
			h.ServeHTTP(w, r)
			return
		}
		sendError(w, http.StatusForbidden, fmt.Errorf("unauthorized"))
		return
	})
}

// timeLogger logs response times
func timeLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		defer log.Printf("It took %s to respond to TimeZone request", time.Since(startTime))
		h.ServeHTTP(res, req)
	})
}

func timeZonesHandler(res http.ResponseWriter, req *http.Request) {
	result, err := sheet.GetTimeZones()
	if err != nil {
		log.Printf("Failed to get time zones: %v\n", err)
		sendError(res, http.StatusBadRequest, fmt.Errorf("Failed to get timezones: %v", err))
		return
	}

	body, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Printf("Failed to marshal json: %v\n", err)
		sendError(res, http.StatusBadRequest, fmt.Errorf("Failed to marshal json: %v\n", err))
		return
	}

	res.Write(body)
}

func usersHandler(res http.ResponseWriter, req *http.Request) {
	result, err := sheet.GetUsers()
	if err != nil {
		log.Printf("Failed to get Users: %v\n", err)
		sendError(res, http.StatusBadRequest, fmt.Errorf("Failed to get users: %v", err))
		return
	}

	body, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal json: %v\n", err)
		sendError(res, http.StatusBadRequest, fmt.Errorf("Failed to marshal json: %v\n", err))
		return
	}

	res.Write(body)
}

func sendError(res http.ResponseWriter, status int, err error) {
	res.WriteHeader(status)
	res.Write([]byte(err.Error()))
}
