package main

import "fmt"

type controller struct {
	IPAddress string `json:"IPAddress"`
	Port      string `json:"Port"`
	Route     string `json:"Route"`
}

type beaconConfig struct {
	PrimaryController   controller `json:"PrimaryController"`
	SecondaryController controller `json:"SecondaryController"`
	CallbackInterval    int        `json:"CallbackInterval"`
	SkewPercentage      int        `json:"SkewPercentage"`
	MaxFailures         int        `json:"MaxFailures"`
}

func (config *beaconConfig) print() {
	fmt.Println("Beacon Config:")
	debugPrintBanner()
	fmt.Printf("Primary Controller: http://%v:%v/%v\n", config.PrimaryController.IPAddress, config.PrimaryController.Port, config.PrimaryController.Route)
	fmt.Printf("Secondary Controller: http://%v:%v/%v\n", config.SecondaryController.IPAddress, config.SecondaryController.Port, config.SecondaryController.Route)
	fmt.Printf("Callback Interval: %v\n", config.CallbackInterval)
	fmt.Printf("Callback Skew: %v\n", config.SkewPercentage)
	fmt.Printf("Maximum Failures: %v\n", config.MaxFailures)
	fmt.Println()
}
