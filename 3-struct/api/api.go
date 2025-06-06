package api

import (
	"bin/bins"
	"bin/config"
	"bin/file"
	"bin/storage"
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

	body := map[string]any{
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

func UpdateBin(id string, filename string, cfg *config.Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error while reading file:  %w", err)
	}

	body := map[string]json.RawMessage{
		"record": data,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Error while marshaling JSON:  %w", err)
	}

	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", id)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("Error while creating request:  %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", cfg.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error while sending request:  %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading response body:  %w", err)
	}

	fmt.Println("Updated bin")
	fmt.Println(string(respData))

	return nil
}

func DeleteBin(id string, cfg *config.Config, store storage.Storage) {
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

	const path = "bins.json"
	storage := &storage.FileStorage{}
	binList, err := storage.LoadBinList(path)
	if err != nil {
		fmt.Println("Error loading bins.json:", err)
		return
	}

	newList := make(bins.BinList, 0)
	for _, b := range binList {
		if b.Id != id {
			newList = append(newList, b)
		}
	}

	err = storage.SaveBinList(newList, path)
	if err != nil {
		fmt.Println("Error saving bins.json:", err)
	}
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

func CreateBinAndReturnID(filePath string, name string, cfg *config.Config) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	var record map[string]any
	err = json.Unmarshal(data, &record)
	if err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	body := map[string]any{
		"record": json.RawMessage(data),
	}
	if name != "" {
		body["name"] = name
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.jsonbin.io/v3/b", strings.NewReader(string(jsonBody)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", cfg.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var result struct {
		Metadata struct {
			ID string `json:"id"`
		} `json:"metadata"`
	}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	bin := bins.NewBin(result.Metadata.ID, true, name)
	if err := bins.Save(bin); err != nil {
		fmt.Println("Warning: failed to save bin locally:", err)
	}

	return result.Metadata.ID, nil
}

func GetBinById(id string, cfg *config.Config) (map[string]any, error) {
	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s/latest", id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("X-Master-Key", cfg.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	var wrapper struct {
		Record map[string]any `json:"record"`
	}
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return wrapper.Record, nil
}
