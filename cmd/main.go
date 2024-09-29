package main

import (
	"CreationDateSaver/config"
	"CreationDateSaver/internal"
	"fmt"
	"log"
)

func main() {
	// Загружаем конфигурацию
	conf, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Error config loading: %v", err)
	}

	// Запускаем отслеживание изменений в папке
	err = internal.WatchFolder(conf.WatchFolder, conf.IncludeSubfolders, conf.MetadataFile)
	if err != nil {
		log.Fatalf("Error folder watching: %v", err)
	}

	fmt.Println("The script is running, changes in the folder are being monitored:", conf.WatchFolder)
}
