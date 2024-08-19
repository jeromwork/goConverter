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

	log.Println("Запуск мониторинга...")
	for {
		config, err := config.LoadConfig(*configPath)
		if err != nil {
			log.Printf("Ошибка загрузки конфигурации: %v", err)
			continue // Если возникла ошибка при загрузке конфигурации, пропускаем итерацию
		}

		monitorFolder.Run(config)
		time.Sleep(10 * time.Second) // Интервал проверки новых файлов
	}
}
