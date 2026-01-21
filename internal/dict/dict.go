package dict

import (
	"embed"
	"strings"
)

//go:embed english10k german10k
var files embed.FS

type Dictionary string

const (
	English10K Dictionary = "english10k"
	German10K  Dictionary = "german10k"
)

func LoadDict(dict Dictionary, top int) []string {
	data, err := files.ReadFile(string(dict))
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.Trim(string(data), "\r\n\t "), "\n")

	return lines[:min(len(lines), top)]
}
