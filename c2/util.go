package main

import (
	"fmt"
	"strings"
)

func debugPrintBanner() {
	banner := strings.Repeat("#", 40)
	fmt.Printf("%v\n", banner)
}
