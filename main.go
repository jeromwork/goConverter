package main

import (
	"flag"
	"log"
	"time"

	"goConverter/config"
	"goConverter/monitorFolder"
)

func main() {
	configPath := flag.String("config", "settings.json", "Путь к файлу конфигурации")
	flag.Parse()

	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	log.Println("Запуск мониторинга...")
	for {
		monitorFolder.Run(config)
		time.Sleep(10 * time.Second) // Интервал проверки новых файлов
	}
}
