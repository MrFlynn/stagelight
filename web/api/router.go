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
	"github.com/gorilla/websocket"
)

var dbhandler *database.DBHandler
var upgrader = websocket.Upgrader{}

func singleDeviceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 8)
	if err != nil {
		http.Error(w, fmt.Sprintf("ID %d is not between 0 and 255", id), http.StatusBadRequest)
		return
	}

	device, err := dbhandler.Get("Devices", uint8(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not find device with ID %d", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(device)

	log.Printf("Sucessfully got device with ID %d", id)
}

func allDeviceHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := dbhandler.GetAll("Devices")
	if err != nil {
		log.Printf("Unable to get list of devices: %s", err)
		http.Error(w, "Coud not get a list of devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(devices)

	log.Println("Sucessfully got a list of all devices")
}

func updateDevices(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Malformed JSON request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = dbhandler.AddMultiple("Devices", body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getVotes(w http.ResponseWriter, r *http.Request) {
	votes, err := dbhandler.Get("Votes", nil)
	if err != nil {
		log.Printf("Unable to get list of votes: %s", err)
		http.Error(w, "Could not get votes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(votes)

	log.Println("Sucessfully got list of votes")
}

func getColors(w http.ResponseWriter, r *http.Request) {
	colors, err := dbhandler.GetAll("Colors")
	if err != nil {
		log.Printf("Unable to get list of colors: %s", err)
		http.Error(w, "Could not get colors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(colors)

	log.Println("Successfully got list of colors")
}

func updateColors(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to parse request body with err: %s", err)
		http.Error(w, "Malformed JSON request", http.StatusBadRequest)

		return
	}

	err = dbhandler.AddMultiple("Colors", body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		log.Printf("Controller unable to update colors: %s", err)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func writeWs(ws *websocket.Conn) {
	d, err := dbhandler.GetAll("Devices")
	if err != nil {
		log.Println("Could not get list of devices to send over websocket")
	}

	ws.WriteJSON(d)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection from %s to websocket", r.RemoteAddr)
		return
	}

	go func() {
		writeWs(ws)

		for range time.Tick(30 * time.Second) {
			writeWs(ws)
		}
	}()
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

	router.HandleFunc("/votes", getVotes).Methods(http.MethodGet)

	router.HandleFunc("/colors", getColors).Methods(http.MethodGet)
	router.HandleFunc("/colors", updateColors).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json;charset=utf-8")

	router.HandleFunc("/ws", wsHandler).Methods(http.MethodGet)

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
