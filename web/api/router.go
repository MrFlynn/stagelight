package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MrFlynn/stagelight/web/api/database"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var dbhandler *database.DBHandler

func singleDeviceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 8)
	if err != nil {
		http.Error(w, fmt.Sprintf("ID %d is not between 0 and 255", id), http.StatusBadRequest)
		return
	}

	device, err := dbhandler.GetDevice(uint8(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not find device with ID %d", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(device)

	log.Printf("Sucessfully got device with ID %d", id)
}

func allDeviceHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := dbhandler.GetAllDevices()
	if err != nil {
		log.Println("Unable to get list of devices.")
		http.Error(w, "Coud not get a list of devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(devices)

	log.Println("Sucessfully got a list of all devices.")
}

func updateDevices(w http.ResponseWriter, r *http.Request) {
	var devices []database.Device

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Malformed JSON request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(body, &devices)
	if err != nil {
		http.Error(w, "Could not process JSON request. Please check format", http.StatusBadRequest)
		return
	}

	err = dbhandler.UpdateDevices(devices)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func createServer() *http.Server {
	log.Println("Creating new database connection...")

	h := database.New("./example.db")
	dbhandler = &h

	log.Println("Creating new router...")

	router := mux.NewRouter()
	router.HandleFunc("/device/{id:[0-9]+}", singleDeviceHandler).Methods(http.MethodGet)
	router.HandleFunc("/device/all", allDeviceHandler).Methods(http.MethodGet)
	router.HandleFunc("/device/update", updateDevices).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json;charset=utf-8")

	// CORS settings.
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodOptions})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"})

	srv := &http.Server{
		Handler:      handlers.CORS(origins, methods, headers)(router),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return srv
}
