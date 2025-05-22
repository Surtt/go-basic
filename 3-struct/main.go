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
	// create := flag.String("create", "", "Create new bin")
	// fileToRead := flag.String("file", "", "File to read")
	// binName := flag.String("name", "", "Bin name")
	binId := flag.String("id", "", "Bin id")
	// list := flag.String("list", "", "Bin list")

	flag.Parse()

	if *get && *binId != "" {
		cfg := config.NewConfig()
		api.GetBinsById(*binId, cfg)
		return
	}

	cfg := config.NewConfig()
	api.GetKey(cfg)

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
