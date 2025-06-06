package bins

import (
	"bin/file"
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Bin struct {
	Id        string    `json:"id"`
	Private   bool      `json:"private"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}

type BinList []Bin

func (list *BinList) Add(bin *Bin) {
	*list = append(*list, *bin)
}

func NewBin(id string, private bool, name string) *Bin {
	bin := &Bin{
		Id:        id,
		Private:   private,
		Name:      name,
		CreatedAt: time.Now(),
	}

	return bin
}

const path = "bins.json"

func Save(bin *Bin) error {
	var list BinList

	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &list)
	}

	list.Add(bin)

	out, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}

func Delete(id string) error {
	const path = "bins.json"

	data, err := file.ReadFile(path)
	if err != nil {
		return err
	}

	var binList BinList
	if err := json.Unmarshal(data, &binList); err != nil {
		return err
	}

	var updated BinList
	found := false
	for _, bin := range binList {
		if bin.Id != id {
			updated = append(updated, bin)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("bin not found in local storage")
	}

	out, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		return err
	}

	return file.WriteFile(path, out)
}
