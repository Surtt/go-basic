package main

import (
	"fmt"
	"time"
)

type Bin struct {
	id        string
	private   bool
	createdAt time.Time
	name      string
}

func newBin(id string, private bool, name string) *Bin {
	bin := &Bin{
		id:        id,
		private:   private,
		name:      name,
		createdAt: time.Now(),
	}

	return bin
}

type BinList []Bin

func (list *BinList) add(bin Bin) {
	*list = append(*list, bin)
}

func main() {
	list := BinList{}
	b := newBin("1234-abcd", true, "Test Bin")
	list.add(*b)
	fmt.Println(list)
}
