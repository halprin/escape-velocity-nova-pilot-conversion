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
	fmt.Printf("Converting %s\n", originalPilotPath)

	resourceForkParser, err := resourcefork.NewParser(originalPilotPath)
	if err != nil {
		fmt.Printf("Error opening resource fork of pilot file: %s", err)
		os.Exit(2)
	}

	primaryPilotResource := resourceForkParser.GetResource("NpïL", 128)
	secondaryPilotResource := resourceForkParser.GetResource("NpïL", 129)

	err = evn.ConvertPilot(originalPilotPath, primaryPilotResource.Data, secondaryPilotResource.Data, secondaryPilotResource.Name)
	if err != nil {
		fmt.Printf("Error converting pilot: %s", err)
		os.Exit(3)
	}
}
