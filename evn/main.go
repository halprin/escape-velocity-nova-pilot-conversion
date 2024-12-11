package evn

import (
	"fmt"
	"os"
	"path/filepath"
)

func ConvertPilot(originalPilotDataPath, originalAltPilotDataPath string) error {
	convertedPilotPath := filepath.Join(filepath.Dir(originalPilotDataPath), filepath.Base(originalPilotDataPath)+".converted.plt")

	convertedPilot, err := os.Create(convertedPilotPath)
	if err != nil {
		return fmt.Errorf("error creating new converted pilot file: %w", err)
	}
	defer convertedPilot.Close()

	err = convertPilotData(originalPilotDataPath, convertedPilot)
	if err != nil {
		return fmt.Errorf("error converting pilot: %w", err)
	}

	err = convertAltPilotData(originalAltPilotDataPath, convertedPilot)
	if err != nil {
		return fmt.Errorf("error converting alt pilot: %w", err)
	}

	err = writeShipName(convertedPilot)
	if err != nil {
		return fmt.Errorf("error writing ship name: %w", err)
	}

	return nil
}

func writeShipName(convertedPilot *os.File) error {
	_, err := convertedPilot.Write([]byte("@PAK\x00"))
	return err
}
