package api

import (
	"bin/config"
	"fmt"
)

func GetKey(cfg *config.Config) {
	fmt.Println(cfg.Key)
}
