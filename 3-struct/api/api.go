package api

import (
	"bin/config"
	"fmt"
	"io"
	"net/http"
)

func GetKey(cfg *config.Config) string {
	return cfg.Key
}

func GetBinsById(id string, cfg *config.Config) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", id), nil)
	if err != nil {
		fmt.Println("Error while creating request:", err)
		return
	}

	req.Header.Set("X-Master-Key", GetKey(cfg))
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error while sending request:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))
}
