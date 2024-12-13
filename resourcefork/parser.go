package resourcefork

import (
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"maps"
	"slices"
)

type Resource struct {
	ID    uint16
	Name  string
	Data  []byte
	Attrs byte
}

type Parser struct {
	data      []byte
	resources map[string]map[uint16]*Resource
}

func NewParser(filePath string) (*Parser, error) {

	resourceData, err := open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing resource fork: %w", err)
	}

	parser := &Parser{
		data:      resourceData,
		resources: make(map[string]map[uint16]*Resource),
	}
	if err := parser.parse(); err != nil {
		return nil, err
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

	// Parse resource map
	mapData := p.data[mapOffset:]
	if len(mapData) < 10 {
		return errors.New("map data too short")
	}

	typeListOffset := 28
	nameListOffset := binary.BigEndian.Uint16(mapData[8:10])

	// Get type list
	if len(mapData) < int(typeListOffset)+2 {
		return errors.New("invalid type list offset")
	}

	macRomanDecoder := charmap.Macintosh.NewDecoder()

	typeCount := int(binary.BigEndian.Uint16(mapData[typeListOffset:typeListOffset+2])) + 1
	currentPos := int(typeListOffset) + 2

	// Parse each resource type
	for typeIndex := 0; typeIndex < typeCount; typeIndex++ {
		if len(mapData) < currentPos+8 {
			return errors.New("invalid resource type data")
		}

		resourceTypeUtf8Bytes, err := macRomanDecoder.Bytes(mapData[currentPos : currentPos+4])
		if err != nil {
			return err
		}
		resourceType := string(resourceTypeUtf8Bytes)

		resourceCount := int(binary.BigEndian.Uint16(mapData[currentPos+4:currentPos+6])) + 1
		refListOffset := binary.BigEndian.Uint16(mapData[currentPos+6 : currentPos+8])

		p.resources[resourceType] = make(map[uint16]*Resource)

		// Parse resources of this type
		refPos := typeListOffset + int(refListOffset)
		for resourceIndex := 0; resourceIndex < resourceCount; resourceIndex++ {
			if len(mapData) < refPos+12 {
				return errors.New("invalid resource reference data")
			}

			resID := binary.BigEndian.Uint16(mapData[refPos : refPos+2])
			nameOffset := binary.BigEndian.Uint16(mapData[refPos+2 : refPos+4])
			attrs := mapData[refPos+4]
			resDataOffset := binary.BigEndian.Uint32(mapData[refPos+4:refPos+8]) & 0x00FFFFFF

			// Get resource name
			var name string
			if nameOffset != 0xFFFF { // 0xFFFF indicates no name
				namePos := typeListOffset + int(nameListOffset) + int(nameOffset)
				if len(mapData) > namePos {
					nameLen := int(mapData[namePos])
					if len(mapData) >= namePos+1+nameLen {

						utf8Bytes, err := macRomanDecoder.Bytes(mapData[namePos+1 : namePos+1+nameLen])
						if err != nil {
							return err
						}
						name = string(utf8Bytes)
					}
				}
			}

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

			p.resources[resourceType][resID] = &Resource{
				ID:    resID,
				Name:  name,
				Data:  resData,
				Attrs: attrs,
			}

			refPos += 12
		}

		currentPos += 8
	}

	return nil
}

// GetResource returns the resource for the given type and ID
func (p *Parser) GetResource(resourceType string, resourceId uint16) *Resource {

	typeMap, ok := p.resources[resourceType]
	if !ok {
		return nil
	}

	resource, ok := typeMap[resourceId]
	if !ok {
		return nil
	}

	return resource
}

// GetTypes returns all resource types
func (p *Parser) GetTypes() []string {
	return slices.Collect(maps.Keys(p.resources))
}

// GetIDs returns all resource IDs for a given type
func (p *Parser) GetIDs(resourceType string) []uint16 {

	typeMap, ok := p.resources[resourceType]
	if !ok {
		return nil
	}

	return slices.Collect(maps.Keys(typeMap))
}
