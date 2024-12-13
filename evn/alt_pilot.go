package evn

import (
	"encoding/binary"
	"fmt"
	"github.com/halprin/escape-velocity-nova-pilot-conversion/endian"
	"os"
)

func convertSecondaryPilotData(secondaryPilotData []byte, convertedPilot *os.File) error {
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(secondaryPilotData)))
	_, err := convertedPilot.Write(lengthBytes)
	if err != nil {
		return fmt.Errorf("error writing alt pilot data size: %w", err)
	}

	decryptedAltPilotData := SimpleDecrypt(secondaryPilotData)
	flippedDecryptedAltPilotData := endian.Flip(decryptedAltPilotData)
	_, err = convertedPilot.Write(flippedDecryptedAltPilotData)
	if err != nil {
		return fmt.Errorf("error writing alt pilot data: %w", err)
	}

	return nil
}
