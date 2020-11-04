package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type beaconResponse struct {
	ID           string       `json:"ID"`
	Beacon       beaconConfig `json:"Beacon"`
	Sunset       sunset       `json:"Sunset"`
	TasksWaiting bool         `json:"TasksWaiting"`
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This service is temporarily unavailable")
}

func viewClients(w http.ResponseWriter, r *http.Request) {
	for _, c := range registeredClients {
		json.NewEncoder(w).Encode(c)
	}
}

func viewClientById(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	for _, c := range registeredClients {
		if c.ID == clientId {
			json.NewEncoder(w).Encode(c)
		}
	}
}

func viewClientTasks(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	for _, c := range registeredClients {
		if c.ID == clientId {
			json.NewEncoder(w).Encode(c.Tasks)
		}
	}
}

func viewClientTasksNew(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	var filteredTasks tasks
	for _, c := range registeredClients {
		if c.ID == clientId {
			filteredTasks = c.Tasks.filterStatus("Created")
		}
	}
	if filteredTasks != nil {
		json.NewEncoder(w).Encode(filteredTasks)
	}
}

func viewClientTasksTasked(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	var filteredTasks tasks
	for _, c := range registeredClients {
		if c.ID == clientId {
			filteredTasks = c.Tasks.filterStatus("Tasked")
		}
	}
	if filteredTasks != nil {
		json.NewEncoder(w).Encode(filteredTasks)
	}
}

func viewClientTasksCompleted(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	var filteredTasks tasks
	for _, c := range registeredClients {
		if c.ID == clientId {
			filteredTasks = c.Tasks.filterStatus("Completed")
		}
	}
	if filteredTasks != nil {
		json.NewEncoder(w).Encode(filteredTasks)
	}
}

func viewClientTasksFailed(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	var filteredTasks tasks
	for _, c := range registeredClients {
		if c.ID == clientId {
			filteredTasks = c.Tasks.filterStatus("Failed")
		}
	}
	if filteredTasks != nil {
		json.NewEncoder(w).Encode(filteredTasks)
	}
}

func viewClientTaskById(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	taskId := mux.Vars(r)["taskId"]
	for _, c := range registeredClients {
		if c.ID == clientId {
			for _, t := range c.Tasks {
				if t.ID == taskId {
					json.NewEncoder(w).Encode(t)
				}
			}
		}
	}
}

func beaconClientRegister(w http.ResponseWriter, r *http.Request) {
	c := registeredClients.addNewClient()

	// just debug defaults for now
	c.BeaconConfig.PrimaryController.IPAddress = "127.0.0.1"
	c.BeaconConfig.PrimaryController.Port = "8080"
	c.BeaconConfig.PrimaryController.Route = "beacon"

	c.BeaconConfig.SecondaryController.IPAddress = "127.0.0.1"
	c.BeaconConfig.SecondaryController.Port = "8080"
	c.BeaconConfig.SecondaryController.Route = "beacon"

	c.BeaconConfig.MaxFailures = 3
	c.BeaconConfig.CallbackInterval = 1
	c.BeaconConfig.SkewPercentage = 0

	c.Sunset.Relative = "5 days"

	c.Tasks.addNewTask("Configure", "")
	c.Tasks.addNewTask("ShellExecute", "ipconfig.exe /all")

	response := beaconResponse{
		ID:           c.ID,
		Beacon:       c.BeaconConfig,
		Sunset:       c.Sunset,
		TasksWaiting: true,
	}

	json.NewEncoder(w).Encode(response)
}

func beaconClientGetConfig(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	c := registeredClients.filterId(clientId)
	if c != nil {
		response := beaconResponse{
			ID:           c.ID,
			Beacon:       c.BeaconConfig,
			Sunset:       c.Sunset,
			TasksWaiting: c.Tasks.filterStatus("Created") != nil,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func beaconClientCheckIn(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	c := registeredClients.filterId(clientId)
	if c != nil {
		// check in

	}
}

func beaconClientPollTasks(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	c := registeredClients.filterId(clientId)
	if c != nil {
		filteredTasks := c.Tasks.filterStatus("Created")
		if filteredTasks != nil {
			json.NewEncoder(w).Encode(filteredTasks)
			filteredTasks.setStatus("Tasked")
		}
	}
}

func beaconClientSubmitTaskById(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]
	taskId := mux.Vars(r)["taskId"]
	c := registeredClients.filterId(clientId)
	if c != nil {
		t := c.Tasks.filterId(taskId)
		if t != nil {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnsupportedMediaType)
			} else {
				json.Unmarshal(body, t)
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}
