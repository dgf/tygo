package dict

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadFile(name string) ([]string, error) {
	data, err := os.ReadFile(filepath.Clean(name))
	if err != nil {
		return nil, fmt.Errorf("read file failed: %w", err)
	}

	var v struct {
		Words []string `json:"words"`
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, fmt.Errorf("unmarshal JSON failed: %w", err)
	}

	return v.Words, nil
}
