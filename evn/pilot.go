package evn

import (
	"encoding/binary"
	"fmt"
	"github.com/halprin/escape-velocity-nova-pilot-conversion/endian"
	"os"
)

func convertPilotData(originalPilotDataPath string, convertedPilot *os.File) error {
	pilotData, err := os.ReadFile(originalPilotDataPath)
	if err != nil {
		return fmt.Errorf("error reading pilot data: %w", err)
	}

	decryptedPilotData := SimpleDecrypt(pilotData)

	firstChunk := readPilotData(&decryptedPilotData, 0x295e)
	flippedDecryptedFirstChunk := endian.Flip(firstChunk)

	flippedDecryptedMissionData := convertMissionData(&decryptedPilotData)

	lastChunk := decryptedPilotData
	flippedDecryptedLastChunk := endian.Flip(lastChunk)

	flippedPilotData := append(flippedDecryptedFirstChunk, flippedDecryptedMissionData...)
	flippedPilotData = append(flippedPilotData, flippedDecryptedLastChunk...)

	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(flippedPilotData)))

	_, err = convertedPilot.Write(lengthBytes)
	if err != nil {
		return fmt.Errorf("error writing pilot data size: %w", err)
	}

	_, err = convertedPilot.Write(flippedPilotData)
	if err != nil {
		return fmt.Errorf("error writing pilot data: %w", err)
	}

	return nil
}

func convertMissionData(decryptedPilotData *[]byte) []byte {
	flippedMissionData := make([]byte, 0)

	for i := 0; i < 16; i++ {
		currentMissionData := make([]byte, 0)

		firstChunk := readPilotData(decryptedPilotData, 0x297e-0x295e)
		currentMissionData = append(currentMissionData, firstChunk...)

		readPilotData(decryptedPilotData, 2) // skip 0x2980 short

		secondChunk := readPilotData(decryptedPilotData, 0x2995-0x2982)
		currentMissionData = append(currentMissionData, secondChunk...)

		readPilotData(decryptedPilotData, 1) // skip 0x2995 Byte

		thirdChunk := readPilotData(decryptedPilotData, 0x3247-0x2996)
		currentMissionData = append(currentMissionData, thirdChunk...)

		flippedCurrentMissionData := endian.Flip(currentMissionData)

		readPilotData(decryptedPilotData, 3) // skip 0x3247 Bytes[3]

		flippedMissionData = append(flippedMissionData, flippedCurrentMissionData...)
	}

	return flippedMissionData
}

func readPilotData(pilotData *[]byte, length int) []byte {
	theData := (*pilotData)[:length]
	*pilotData = (*pilotData)[length:]
	return theData
}
