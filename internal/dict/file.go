package dict

import (
	"encoding/json"
	"os"
)

func LoadFile(fileName string) []string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	type languageFile struct {
		Words []string `json:"words"`
	}

	var lf languageFile

	err = json.Unmarshal(data, &lf)
	if err != nil {
		panic(err)
	}

	return lf.Words
}
