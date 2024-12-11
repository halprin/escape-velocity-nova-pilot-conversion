package main

import (
	"fmt"
	"github.com/halprin/escape-velocity-nova-pilot-conversion/evn"
	"os"
)


func main() {
	if len(os.Args) != 3 {
		fmt.Println("Provide a path to the pilot data and alt player data file")
		os.Exit(1)
	}

	originalPilotDataPath := os.Args[1]
	originalAltPilotDataPath := os.Args[2]
	fmt.Printf("Converting %s and %s\n", originalPilotDataPath, originalAltPilotDataPath)

	err := evn.ConvertPilot(originalPilotDataPath, originalAltPilotDataPath)
	if err != nil {
		fmt.Printf("Error converting pilot: %s", err)
		os.Exit(2)
	}
}
