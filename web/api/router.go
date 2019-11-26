package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MrFlynn/stagelight/web/api/database"
	"github.com/gorilla/mux"
)

func startAndServe() {
	handler := database.New("./example.db")

	router := mux.NewRouter()
	router.HandleFunc("/device/{id: [0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 8)
		if err != nil {
			http.Error(w, fmt.Sprintf("ID %d is not between 0 and 255", id), http.StatusBadRequest)
			return
		}

		device, err := handler.GetDevice(uint8(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not find device with ID %d", id), http.StatusNotFound)
		}

		json.NewEncoder(w).Encode(device)
	}).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
