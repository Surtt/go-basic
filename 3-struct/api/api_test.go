package api_test

import (
	"bin/api"
	"bin/config"
	"encoding/json"
	"os"
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
	data := map[string]any{"foo": "bar"}
	writeTempJSON(t, data, "create.json")
	defer os.Remove("create.json")

	id, err := api.CreateBinAndReturnID("create.json", "TestBin", cfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer api.DeleteBin(id, cfg)

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
	original := map[string]any{"key": "old"}
	updated := map[string]any{"key": "new"}

	writeTempJSON(t, original, "original.json")
	defer os.Remove("original.json")

	id, err := api.CreateBinAndReturnID("original.json", "UpdateTest", cfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer api.DeleteBin(id, cfg)

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
	data := map[string]any{"toDelete": true}
	writeTempJSON(t, data, "del.json")
	defer os.Remove("del.json")

	id, err := api.CreateBinAndReturnID("del.json", "DelTest", cfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	api.DeleteBin(id, cfg)

	record, err := api.GetBinById(id, cfg)
	if err == nil && record != nil {
		t.Errorf("Expected bin to be deleted, but got: %v", record)
	}
}
