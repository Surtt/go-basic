package main

import (
	"bin/bins"
	"fmt"
)

func main() {
	list := bins.BinList{}
	b := bins.NewBin("1234-abcd", true, "Test Bin")
	list.Add(b)
	fmt.Println(list)
}
