package main

import "fmt"

type sunset struct {
	Absolute string `json:"Absolute"`
	Relative string `json:"Relative"`
}

func (s *sunset) print() {
	fmt.Println("Sunset:")
	debugPrintBanner()
	fmt.Printf("Absolute: %v\n", s.Absolute)
	fmt.Printf("Relative: %v\n", s.Relative)
	fmt.Println()
}
