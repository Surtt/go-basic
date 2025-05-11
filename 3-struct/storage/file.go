package storage

import (
	"bin/bins"
	"encoding/json"
	"errors"
	"os"

	"bin/file"
)

type FileStorage struct{}

func (fs *FileStorage) SaveBinList(binList bins.BinList, path string) error {
	data, err := json.Marshal(binList)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (fs *FileStorage) LoadBinList(path string) (bins.BinList, error) {
	if !file.IsJSON(path) {
		return nil, errors.New("file is not json")
	}
	data, err := file.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var binList bins.BinList
	err = json.Unmarshal(data, &binList)
	if err != nil {
		return nil, err
	}
	return binList, nil
}
