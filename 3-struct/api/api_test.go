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
	bytes, _ := json.Marshal(data)
	_ = os.WriteFile(filename, bytes, 0644)
}

func TestCreateBinAndGet(t *testing.T) {
	cfg := config.NewConfig()
	data := map[string]any{"foo": "bar"}
	writeTempJSON(t, data, "temp-create.json")
	defer os.Remove("temp-create.json")

	id := api.CreateBinAndReturnID("temp-create.json", "TestBin", cfg)
	if id == "" {
		t.Fatal("Bin was not created")
	}

	res := api.GetBinById(id, cfg)

	inner, ok := res["record"].(map[string]any)
	if !ok || inner["foo"] != "bar" {
		t.Errorf("Expected foo=bar, got %v", res)
	}

	api.DeleteBin(id, cfg)
}

func TestUpdateBinAndGet(t *testing.T) {
	cfg := config.NewConfig()
	initial := map[string]any{"key": "old"}
	updated := map[string]any{"key": "new"}

	writeTempJSON(t, initial, "temp-init.json")
	defer os.Remove("temp-init.json")

	id := api.CreateBinAndReturnID("temp-init.json", "UpdateTest", cfg)
	defer api.DeleteBin(id, cfg)

	writeTempJSON(t, updated, "temp-update.json")
	defer os.Remove("temp-update.json")

	api.UpdateBin(id, "temp-update.json", cfg)
	res := api.GetBinById(id, cfg)

	inner, ok := res["record"].(map[string]any)
	if !ok || inner["key"] != "new" {
		t.Errorf("Expected key=new, got %v", res)
	}
}

func TestDeleteBin(t *testing.T) {
	cfg := config.NewConfig()
	data := map[string]any{"delete": true}
	writeTempJSON(t, data, "temp-delete.json")
	defer os.Remove("temp-delete.json")

	id := api.CreateBinAndReturnID("temp-delete.json", "DeleteTest", cfg)
	api.DeleteBin(id, cfg)

	result := api.GetBinById(id, cfg)
	if result != nil {
		t.Error("Bin was not deleted")
	}
}
