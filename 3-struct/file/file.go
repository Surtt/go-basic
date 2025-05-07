package file

import (
	"os"
	"strings"
)

func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func IsJSON(path string) bool {
	return strings.HasSuffix(path, ".json")
}
