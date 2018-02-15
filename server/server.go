package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

	s.HandleFunc("/timezones", TimeZonesHandler)
	s.HandleFunc("/users", UsersHandler)

	var err error
	sheet, err = initSheet()
	if err != nil {
		return err
	}

	log.Println("Starting server...")
	return http.ListenAndServe(":3000", s)
}

func TimeZonesHandler(res http.ResponseWriter, req *http.Request) {

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

func UsersHandler(res http.ResponseWriter, req *http.Request) {
	result, err := sheet.GetUsers()
	if err != nil {
		log.Printf("Failed to get Users: %v\n", err)
		sendError(res, http.StatusBadRequest, fmt.Errorf("Failed to get users: %v", err))
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

func sendError(res http.ResponseWriter, status int, err error) {
	res.WriteHeader(status)
	res.Write([]byte(err.Error()))
}
