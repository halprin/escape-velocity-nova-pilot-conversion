package resourcefork

import (
	"fmt"
	"os"
	"path/filepath"
)

func open(filePath string) ([]byte, error) {
	rsrcPath := filepath.Join(filePath, "..namedfork/rsrc")

	data, err := os.ReadFile(rsrcPath)
	if err != nil {
		return nil, fmt.Errorf("error reading resource fork: %w", err)
	}

	return data, nil
}
