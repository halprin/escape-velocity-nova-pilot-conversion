package main

import (
	"fmt"
	"github.com/halprin/escape-velocity-nova-pilot-conversion/evn"
	"github.com/halprin/escape-velocity-nova-pilot-conversion/resourcefork"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Provide a path to the pilot data and alt player data file")
		os.Exit(1)
	}

	originalPilotPath := os.Args[1]
	fmt.Printf("Converting pilot %s\n", originalPilotPath)

	resourceForkParser, err := resourcefork.NewParser(originalPilotPath)
	if err != nil {
		fmt.Printf("Error opening resource fork of pilot file: %s", err)
		os.Exit(2)
	}

	primaryPilotResource := resourceForkParser.GetResource("NpïL", 128)
	secondaryPilotResource := resourceForkParser.GetResource("NpïL", 129)

	convertedPilotPath := originalPilotPath + ".converted.plt"

	convertedPilot, err := os.Create(convertedPilotPath)
	if err != nil {
		fmt.Printf("Error creating new converted pilot file: %s", err)
		os.Exit(3)
	}
	defer convertedPilot.Close()

	err = evn.ConvertPilot(convertedPilot, primaryPilotResource.Data, secondaryPilotResource.Data, secondaryPilotResource.Name)
	if err != nil {
		fmt.Printf("Error converting pilot: %s", err)
		os.Exit(4)
	}

	fmt.Printf("Converted pilot to %s\n", convertedPilotPath)
}
