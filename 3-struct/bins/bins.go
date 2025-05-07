package bins

import "time"

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
