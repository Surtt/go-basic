package main

import (
	"bin/api"
	"bin/bins"
	"bin/config"
	"bin/file"
	"bin/storage"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Не удалось загрузить переменные окружения")
	}
	get := flag.Bool("get", false, "Get bin by id")
	create := flag.Bool("create", false, "Create bin from file")
	update := flag.Bool("update", false, "Update bin by ID")
	delete := flag.Bool("delete", false, "Delete bin by ID")
	list := flag.Bool("list", false, "List all bins")

	fileToRead := flag.String("file", "", "Path to JSON file")
	binName := flag.String("name", "", "Name of bin")
	binId := flag.String("id", "", "Bin id")

	flag.Parse()

	cfg := config.NewConfig()

	if *get && *binId != "" {
		record, err := api.GetBinById(*binId, cfg)
		if err != nil {
			fmt.Println("Error getting bin:", err)
			return
		}
		data, _ := json.MarshalIndent(record, "", "  ")
		fmt.Println(string(data))
		return
	}

	if *create && *fileToRead != "" {
		id, err := api.CreateBinAndReturnID(*fileToRead, *binName, cfg)
		if err != nil {
			fmt.Println("❌ Error while creating bin:", err)
			return
		}

		bin := bins.NewBin(id, true, *binName)
		if err := bins.Save(bin); err != nil {
			fmt.Println("❌ Error while saving bin to bins.json:", err)
			return
		}

		fmt.Println("✅ Bin created with ID:", id)
		return
	}

	if *update && *fileToRead != "" && *binId != "" {
		api.UpdateBin(*binId, *fileToRead, cfg)
		return
	}

	if *delete && *binId != "" {
		store := &storage.FileStorage{}
		api.DeleteBin(*binId, cfg, store)
		return
	}

	if *list {
		api.ListBins(cfg)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Укажите имя файла с бин, например: newbin.json")
		return
	}

	newBinFile := os.Args[1]

	if !file.IsJSON(newBinFile) {
		fmt.Println("Файл нового Bin должен быть JSON.")
		return
	}

	data, err := file.ReadFile(newBinFile)
	if err != nil {
		fmt.Println("Не удалось прочитвть файл", err)
		return
	}

	var newBins bins.BinList
	err = json.Unmarshal(data, &newBins)
	if err != nil {
		fmt.Println("Не удалось преобразовать данные", err)
		return
	}

	fmt.Printf("Прочитано %d новых бин(ов) из %s\n", len(newBins), newBinFile)

	var storage storage.Storage = &storage.FileStorage{}

	binList, err := storage.LoadBinList("bins.json")
	if err != nil {
		fmt.Println("bins.json не найден или пустой. Будет создан новый список.")
		binList = bins.BinList{}
	}

	for _, bin := range newBins {
		binList.Add(&bin)
	}

	err = storage.SaveBinList(binList, "bins.json")
	if err != nil {
		fmt.Println("Ошибка при сохранении bins.json:", err)
		return
	}

	fmt.Printf("Список успешно обновлён. Теперь в bins.json %d бин(ов).\n", len(binList))
}
