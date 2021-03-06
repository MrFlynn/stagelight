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
var deviceUpdateStream = make(chan []database.Device, 10)

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

	// If the previous if statement didn't error, we know this should be error free.
	devicesPayload := []database.Device{}
	json.Unmarshal(body, &devicesPayload)
	deviceUpdateStream <- devicesPayload

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

func updateVotes(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to parse request body with err: %s", err)
		http.Error(w, "Malformed JSON request", http.StatusBadRequest)

		return
	}

	err = dbhandler.AddMultiple("Votes", body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		log.Printf("Controller unable to update votes: %s", err)

		return
	}

	w.WriteHeader(http.StatusCreated)
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

func addNewColor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var objmap map[string]*json.RawMessage
	var name string
	var sequence []uint32
	
	err = json.Unmarshal(body, &objmap)
	err = json.Unmarshal(*objmap["name"], &name)
	err = json.Unmarshal(*objmap["sequence"], &sequence)
	if err != nil {
		log.Printf("Unable to parse request body with err: %s", err)
		http.Error(w, "Malformed JSON request", http.StatusBadRequest)

		return
	}

	// This is super inneficient, but I don't time to implement a better method at the moment.
	colors, err := dbhandler.GetAll("Colors")
	if err != nil {
		log.Printf("Unable to get list of colors: %s", err)
		http.Error(w, "Could not get a list of colors", http.StatusInternalServerError)
		
		return
	}

	newColor := database.Color{
		ID: uint8(len(colors) + 1),
		Name: name,
		Sequence: sequence,
	}
	val, _ := json.Marshal(newColor)
	err = dbhandler.Add("Colors", val)

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		log.Printf("Controller unable to add new color: %s", err)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func wsBridgeHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection from %s to websocket", r.RemoteAddr)
		return
	}

	go func() {
		for {
			// Only send updates as they are available.
			devices := <-deviceUpdateStream
			ws.WriteJSON(dbhandler.CreatePayload(devices))
		}
	}()
}

func wsVoteHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection from %s to websocket", r.RemoteAddr)
		return
	}

	go func() {
		for range time.Tick(2 * time.Second) {
			votes, err := dbhandler.Get("Votes", nil)
			if err != nil {
				log.Printf("Could not get list of votes.")
				continue
			}

			ws.WriteJSON(votes)
		}
	}()
}

func createServer(addr string, port uint, databasePath string) *http.Server {
	log.Println("Creating new database connection...")

	h := database.New(databasePath)
	dbhandler = &h

	log.Println("Creating new router...")

	router := mux.NewRouter()

	router.HandleFunc("/device/{id:[0-9]+}", singleDeviceHandler).Methods(http.MethodGet)
	router.HandleFunc("/device/all", allDeviceHandler).Methods(http.MethodGet)
	router.HandleFunc("/device/update", updateDevices).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json;charset=utf-8")

	router.HandleFunc("/votes", getVotes).Methods(http.MethodGet)
	router.HandleFunc("/votes", updateVotes).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json;charset=utf-8")

	router.HandleFunc("/colors", getColors).Methods(http.MethodGet)
	router.HandleFunc("/colors", updateColors).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json;charset=utf-8")
	router.HandleFunc("/colors/new", addNewColor).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json;charset=utf-8")

	router.HandleFunc("/ws/bridge", wsBridgeHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws/votes", wsVoteHandler).Methods(http.MethodGet)

	// CORS settings.
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodOptions})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"})

	srv := &http.Server{
		Handler:      handlers.CORS(origins, methods, headers)(router),
		Addr:         fmt.Sprintf("%s:%d", addr, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return srv
}
