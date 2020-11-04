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
	router.HandleFunc("/view/client/{clientId}", viewClientById).Methods("GET")
	router.HandleFunc("/view/client/{clientId}/tasks", viewClientTasks).Methods("GET")
	router.HandleFunc("/view/client/{clientId}/tasks/new", viewClientTasksNew).Methods("GET")
	router.HandleFunc("/view/client/{clientId}/tasks/tasked", viewClientTasksTasked).Methods("GET")
	router.HandleFunc("/view/client/{clientId}/tasks/completed", viewClientTasksCompleted).Methods("GET")
	router.HandleFunc("/view/client/{clientId}/tasks/failed", viewClientTasksFailed).Methods("GET")

	router.HandleFunc("/beacon", beaconClientRegister).Methods("POST")
	router.HandleFunc("/beacon/{clientId}", beaconClientGetConfig).Methods("GET")
	router.HandleFunc("/beacon/{clientId}", beaconClientCheckIn).Methods("POST")
	router.HandleFunc("/beacon/{clientId}/tasks", beaconClientPollTasks).Methods("GET")
	router.HandleFunc("/beacon/{clientId}/task/{taskId}", beaconClientSubmitTaskById).Methods("POST")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
