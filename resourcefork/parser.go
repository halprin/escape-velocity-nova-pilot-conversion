package resourcefork

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Parser struct {
	data      []byte
	resources map[string]map[uint16][]byte
}

func NewParser(filePath string) (*Parser, error) {

	resourceData, err := open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing resource fork: %w", err)
	}

	parser := &Parser{
		data:      resourceData,
		resources: make(map[string]map[uint16][]byte),
	}

	err = parser.parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing resource fork: %w", err)
	}

	return parser, nil
}

func (p *Parser) parse() error {
	if len(p.data) < 16 {
		return errors.New("resource fork data too short")
	}

	// Resource fork header
	dataOffset := binary.BigEndian.Uint32(p.data[0:4])
	mapOffset := binary.BigEndian.Uint32(p.data[4:8])
	// dataLength := binary.BigEndian.Uint32(p.data[8:12])
	// mapLength := binary.BigEndian.Uint32(p.data[12:16])

	// Parse resource map
	mapData := p.data[mapOffset:]
	if len(mapData) < 10 {
		return errors.New("map data too short")
	}

	// Skip header bytes in map
	// attributes := binary.BigEndian.Uint16(mapData[4:6])
	typeListOffset := binary.BigEndian.Uint16(mapData[6:8])
	// nameListOffset := binary.BigEndian.Uint16(mapData[8:10])

	// Get type list
	if len(mapData) < int(typeListOffset)+2 {
		return errors.New("invalid type list offset")
	}

	typeCount := int(binary.BigEndian.Uint16(mapData[typeListOffset:typeListOffset+2])) + 1
	currentPos := int(typeListOffset) + 2

	// Parse each resource type
	for i := 0; i < typeCount; i++ {
		if len(mapData) < currentPos+8 {
			return errors.New("invalid resource type data")
		}

		// Get resource type (4 chars) and count
		resType := string(mapData[currentPos : currentPos+4])
		resCount := int(binary.BigEndian.Uint16(mapData[currentPos+4:currentPos+6])) + 1
		refListOffset := binary.BigEndian.Uint16(mapData[currentPos+6 : currentPos+8])

		p.resources[resType] = make(map[uint16][]byte)

		// Parse resources of this type
		refPos := int(typeListOffset) + int(refListOffset)
		for j := 0; j < resCount; j++ {
			if len(mapData) < refPos+12 {
				return errors.New("invalid resource reference data")
			}

			resID := binary.BigEndian.Uint16(mapData[refPos : refPos+2])
			// nameOffset := binary.BigEndian.Uint16(mapData[refPos+2 : refPos+4])
			// resAttributes := mapData[refPos+4]
			resDataOffset := binary.BigEndian.Uint32(mapData[refPos+4:refPos+8]) & 0x00FFFFFF

			// Get resource data
			dataPos := int(dataOffset) + int(resDataOffset)
			if len(p.data) < dataPos+4 {
				return errors.New("invalid resource data offset")
			}

			resDataLen := binary.BigEndian.Uint32(p.data[dataPos : dataPos+4])
			if len(p.data) < dataPos+4+int(resDataLen) {
				return errors.New("invalid resource data length")
			}

			resData := p.data[dataPos+4 : dataPos+4+int(resDataLen)]
			p.resources[resType][resID] = resData

			refPos += 12
		}

		currentPos += 8
	}

	return nil
}

// GetResource returns the resource data for the given type and ID
func (p *Parser) GetResource(resType string, resID uint16) []byte {
	if typeMap, ok := p.resources[resType]; ok {
		if data, ok := typeMap[resID]; ok {
			return data
		}
	}
	return nil
}

// GetTypes returns all resource types
func (p *Parser) GetTypes() []string {
	types := make([]string, 0, len(p.resources))
	for t := range p.resources {
		types = append(types, t)
	}
	return types
}

// GetIDs returns all resource IDs for a given type
func (p *Parser) GetIDs(resType string) []uint16 {
	if typeMap, ok := p.resources[resType]; ok {
		ids := make([]uint16, 0, len(typeMap))
		for id := range typeMap {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}
