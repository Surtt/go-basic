package storage

import "bin/bins"

type Storage interface {
	SaveBinList(binList bins.BinList, path string) error
	LoadBinList(path string) (bins.BinList, error)
}
