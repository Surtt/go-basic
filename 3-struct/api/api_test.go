package api_test

import (
	"bin/api"
	"bin/bins"
	"bin/config"
	"bin/storage"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("./../.env")

	code := m.Run()
	os.Exit(code)
}

func writeTempJSON(t *testing.T, data map[string]any, filename string) {
	t.Helper()
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		t.Fatalf("write file error: %v", err)
	}
}

func TestCreateBinAndGet(t *testing.T) {
	cfg := config.NewConfig()
	store := &storage.FileStorage{}
	data := map[string]any{"foo": "bar"}
	writeTempJSON(t, data, "create.json")
	defer os.Remove("create.json")

	id, err := api.CreateBinAndReturnID("create.json", "TestBin", cfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer api.DeleteBin(id, cfg, store)

	record, err := api.GetBinById(id, cfg)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if record["foo"] != "bar" {
		t.Errorf("Expected foo=bar, got %v", record)
	}
}

func TestUpdateBinAndGet(t *testing.T) {
	cfg := config.NewConfig()
	store := &storage.FileStorage{}
	original := map[string]any{"key": "old"}
	updated := map[string]any{"key": "new"}

	writeTempJSON(t, original, "original.json")
	defer os.Remove("original.json")

	id, err := api.CreateBinAndReturnID("original.json", "UpdateTest", cfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer api.DeleteBin(id, cfg, store)

	writeTempJSON(t, updated, "updated.json")
	defer os.Remove("updated.json")

	api.UpdateBin(id, "updated.json", cfg)

	record, err := api.GetBinById(id, cfg)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if record["key"] != "new" {
		t.Errorf("Expected key=new, got %v", record)
	}
}

func TestDeleteBinConfirm(t *testing.T) {
	cfg := config.NewConfig()
	store := &storage.FileStorage{}
	data := map[string]any{"toDelete": true}
	writeTempJSON(t, data, "del.json")
	defer os.Remove("del.json")

	id, err := api.CreateBinAndReturnID("del.json", "DelTest", cfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	api.DeleteBin(id, cfg, store)

	record, err := api.GetBinById(id, cfg)
	if err == nil && record != nil {
		t.Errorf("Expected bin to be deleted, but got: %v", record)
	}
}

func TestCreateBin_SavesToLocalFile(t *testing.T) {
	cfg := config.NewConfig()
	store := &storage.FileStorage{}
	tempFile := "temp_create.json"
	binsFile := "bins.json"

	data := map[string]any{"foo": "bar"}
	writeTempJSON(t, data, tempFile)
	defer os.Remove(tempFile)
	defer os.Remove(binsFile)

	id, err := api.CreateBinAndReturnID(tempFile, "LocalSaveTest", cfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer api.DeleteBin(id, cfg, store)

	bin := bins.NewBin(id, true, "LocalSaveTest")
	if err := bins.Save(bin); err != nil {
		t.Fatalf("Save to bins.json failed: %v", err)
	}

	content, err := os.ReadFile(binsFile)
	if err != nil {
		t.Fatalf("Failed to read bins.json: %v", err)
	}
	if !json.Valid(content) {
		t.Error("Invalid JSON in bins.json")
	}
	if !containsID(content, id) {
		t.Errorf("bins.json does not contain created ID %s", id)
	}
}

func TestDeleteBin_RemovesFromLocalFile(t *testing.T) {
	cfg := config.NewConfig()
	store := &storage.FileStorage{}
	tempFile := "temp_delete.json"
	binsFile := "bins.json"

	data := map[string]any{"key": "value"}
	writeTempJSON(t, data, tempFile)
	defer os.Remove(tempFile)
	defer os.Remove(binsFile)

	id, err := api.CreateBinAndReturnID(tempFile, "DeleteLocalTest", cfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	bin := bins.NewBin(id, true, "DeleteLocalTest")
	if err := bins.Save(bin); err != nil {
		t.Fatalf("Failed to save bin: %v", err)
	}

	err = bins.Delete(id)
	if err != nil {
		t.Fatalf("Delete from bins.json failed: %v", err)
	}

	content, err := os.ReadFile(binsFile)
	if err != nil {
		t.Fatalf("Failed to read bins.json: %v", err)
	}
	if containsID(content, id) {
		t.Errorf("bins.json still contains deleted ID %s", id)
	}

	api.DeleteBin(id, cfg, store)
}

func containsID(data []byte, id string) bool {
	return strings.Contains(string(data), id)
}
