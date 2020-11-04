package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var registeredClients = clients{}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/view/clients", viewClients).Methods("GET")
	router.HandleFunc("/view/client/{clientID}", viewClientByID).Methods("GET")
	router.HandleFunc("/view/client/{clientID}/tasks", viewClientTasks).Methods("GET")
	router.HandleFunc("/view/client/{clientID}/tasks/new", viewClientTasksNew).Methods("GET")
	router.HandleFunc("/view/client/{clientID}/tasks/tasked", viewClientTasksTasked).Methods("GET")
	router.HandleFunc("/view/client/{clientID}/tasks/completed", viewClientTasksCompleted).Methods("GET")
	router.HandleFunc("/view/client/{clientID}/tasks/failed", viewClientTasksFailed).Methods("GET")

	router.HandleFunc("/beacon", beaconClientRegister).Methods("POST")
	router.HandleFunc("/beacon/{clientID}", beaconClientGetConfig).Methods("GET")
	router.HandleFunc("/beacon/{clientID}", beaconClientCheckIn).Methods("POST")
	router.HandleFunc("/beacon/{clientID}/tasks", beaconClientPollTasks).Methods("GET")
	router.HandleFunc("/beacon/{clientID}/task/{taskID}", beaconClientSubmitTaskByID).Methods("POST")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
