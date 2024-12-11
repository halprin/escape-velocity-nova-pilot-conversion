package evn

import (
	"encoding/binary"
	"fmt"
	"github.com/halprin/escape-velocity-nova-pilot-conversion/endian"
	"os"
)

func convertAltPilotData(originalAltPilotDataPath string, convertedPilot *os.File) error {
	altPilotData, err := os.ReadFile(originalAltPilotDataPath)
	if err != nil {
		return fmt.Errorf("error reading alt pilot data: %w", err)
	}

	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(altPilotData)))
	_, err = convertedPilot.Write(lengthBytes)
	if err != nil {
		return fmt.Errorf("error writing alt pilot data size: %w", err)
	}

	decryptedAltPilotData := SimpleDecrypt(altPilotData)
	flippedDecryptedAltPilotData := endian.Flip(decryptedAltPilotData)
	_, err = convertedPilot.Write(flippedDecryptedAltPilotData)
	if err != nil {
		return fmt.Errorf("error writing alt pilot data: %w", err)
	}

	return nil
}
