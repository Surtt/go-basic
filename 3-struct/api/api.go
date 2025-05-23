package api

import (
	"bin/bins"
	"bin/config"
	"bin/file"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

func CreateBin(filename string, name string, cfg *config.Config) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error while reading file:", err)
		return
	}

	body := map[string]interface{}{
		"record": json.RawMessage(data),
	}

	if name != "" {
		body["name"] = name
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error while marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.jsonbin.io/v3/b", strings.NewReader(string(jsonBody)))
	if err != nil {
		fmt.Println("Error while creating request:", err)
		return
	}

	req.Header.Set("X-Master-Key", GetKey(cfg))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error while sending request:", err)
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response body:", err)
		return
	}

	fmt.Println("Created new bin")
	fmt.Println(string(respBody))
}

func UpdateBin(id string, filename string, cfg *config.Config) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error while reading file", err)
		return
	}

	body := map[string]json.RawMessage{
		"record": data,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error while marshaling JSON", err)
		return
	}

	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", id)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(jsonBody)))
	if err != nil {
		fmt.Println("Error while creating request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", cfg.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error while sending request", err)
		return
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response body", err)
		return
	}

	fmt.Println("Updated bin")
	fmt.Println(string(respData))
}

func DeleteBin(id string, cfg *config.Config) {
	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println("Error while creating request", err)
		return
	}

	req.Header.Set("X-Master-Key", cfg.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error while sending request", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response body", err)
		return
	}

	fmt.Println("Deleted bin")
	fmt.Println(string(body))
}

func ListBins(cfg *config.Config) {
	const path = "bins.json"

	if !file.IsJSON(path) {
		fmt.Println("File bins.json not found")
		return
	}

	data, err := file.ReadFile(path)
	if err != nil {
		fmt.Println("Error while reading file", err)
		return
	}

	var binList bins.BinList
	if err := json.Unmarshal(data, &binList); err != nil {
		fmt.Println("Error while unmarshaling JSON", err)
		return
	}

	if len(binList) == 0 {
		fmt.Println("No bins found")
		return
	}

	fmt.Println("List of bins:")
	for i, bin := range binList {
		fmt.Printf("%d. %s | id=%s | private=%t | created=%s\n",
			i+1, bin.Name, bin.Id, bin.Private, bin.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
