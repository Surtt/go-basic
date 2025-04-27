package bins

import "time"

type Bin struct {
	id        string
	private   bool
	createdAt time.Time
	name      string
}

type BinList []Bin

func (list *BinList) Add(bin *Bin) {
	*list = append(*list, *bin)
}

func NewBin(id string, private bool, name string) *Bin {
	bin := &Bin{
		id:        id,
		private:   private,
		name:      name,
		createdAt: time.Now(),
	}

	return bin
}
