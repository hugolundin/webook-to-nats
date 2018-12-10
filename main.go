package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/go-nats"
)

var nc *nats.Conn

// Forward messages from HTTP to NATS.
func Forward(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)

	if r.Body == nil {
		http.Error(w, "A body is required.", http.StatusBadRequest)
		return
	}

	bytes, _ := ioutil.ReadAll(r.Body)
	if len(bytes) <= 0 {
		http.Error(w, "Invalid length", http.StatusBadRequest)
		return
	}

	ret := nc.Publish("topic", bytes)
	if ret != nil {
		http.Error(w, "NATS publish failed.", http.StatusInternalServerError)
	}
}

func main() {
	nc, _ = nats.Connect("nats://0.0.0.0:4222")

	router := mux.NewRouter()

	router.PathPrefix("/").HandlerFunc(Forward)

	router.HandleFunc("/{topic}", Forward).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
