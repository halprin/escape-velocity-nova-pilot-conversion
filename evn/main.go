package evn

import (
	"fmt"
	"os"
)

func ConvertPilot(convertedPilot *os.File, primaryPilotData []byte, secondaryPilotData []byte, shipName string) error {
	err := convertPrimaryPilotData(primaryPilotData, convertedPilot)
	if err != nil {
		return fmt.Errorf("error converting primary pilot data: %w", err)
	}

	err = convertSecondaryPilotData(secondaryPilotData, convertedPilot)
	if err != nil {
		return fmt.Errorf("error converting secondary pilot data: %w", err)
	}

	err = writeShipName(shipName, convertedPilot)
	if err != nil {
		return fmt.Errorf("error writing ship name: %w", err)
	}

	return nil
}

func writeShipName(shipName string, convertedPilot *os.File) error {
	_, err := convertedPilot.WriteString(shipName)
	if err != nil {
		return err
	}

	_, err = convertedPilot.Write([]byte("\x00"))
	return err
}
