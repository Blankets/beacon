package main

import "fmt"

type client struct {
	ID           string       `json:"ID"`
	IPAddress    string       `json:"IPAddress"`
	BeaconConfig beaconConfig `json:"Beacon"`
	Sunset       sunset       `json:"Sunset"`
	Tasks        tasks        `json:"Tasks"`
}

type clients []*client

func createClient() *client {
	var newClient client

	newClient.ID = getNewId()
	newClient.IPAddress = ""
	newClient.Sunset = sunset{}
	newClient.BeaconConfig = beaconConfig{}
	newClient.Tasks = tasks{}

	return &newClient
}

func (cs *clients) appendClient(c *client) {
	*cs = append(*cs, c)
}

func (cs *clients) addNewClient() *client {
	newClient := createClient()
	cs.appendClient(newClient)

	return newClient
}

func (cs *clients) filterId(id string) *client {
	for _, c := range *cs {
		if c.ID == id {
			return c
		}
	}
	return nil
}

func (c *client) print() {
	debugPrintBanner()
	fmt.Printf("Client ID:\t%v\n", c.ID)
	fmt.Printf("IP Address:\t%v\n", c.IPAddress)
	fmt.Println()

	c.BeaconConfig.print()
	c.Sunset.print()
	c.Tasks.print()
	debugPrintBanner()
	fmt.Println()
}

func (cs *clients) print() {
	for _, c := range *cs {
		c.print()
	}
}
