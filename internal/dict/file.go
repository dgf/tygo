package dict

import (
	"encoding/json"
	"os"
)

func LoadFile(name string) ([]string, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var lf struct {
		Words []string `json:"words"`
	}

	err = json.Unmarshal(data, &lf)
	if err != nil {
		return nil, err
	}

	return lf.Words, nil
}
