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

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func IsJSON(path string) bool {
	return strings.HasSuffix(path, ".json")
}
