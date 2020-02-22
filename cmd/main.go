package main

import (
	"flag"
	"log"
	"yandex-catalog-handler/internal/concator"
	"yandex-catalog-handler/internal/result"
	"yandex-catalog-handler/pkg/storage"

	"yandex-catalog-handler/internal/consumer"
	"yandex-catalog-handler/internal/loader"
	"yandex-catalog-handler/internal/server"
	"yandex-catalog-handler/pkg/config"
)

func main() {
	configPath := flag.String("config-path", "", "path to config .yaml file")

	flag.Parse()

	if *configPath == "" {
		log.Fatalf("config-path is empty")
		return
	}

	cfg := config.Config{}
	cfg.ReadConfigFromPath(*configPath)

	loaderService := loader.New(cfg)

	concatorService := concator.New(cfg)

	consumerService := consumer.New(cfg, loaderService)

	db, err := storage.New(cfg)

	resultRepo := result.NewRepository(db)
	resultService := result.NewService(resultRepo)

	if err != nil {
		log.Fatalf("Error when try storage.New, err: %s", err)
	}

	s := server.NewServer(cfg, consumerService, loaderService, concatorService, resultService)

	err = s.Run()

	if err != nil {
		log.Fatalf("Error when try server.NewServer, err: %s", err)
	}
}
