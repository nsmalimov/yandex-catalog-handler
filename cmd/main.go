package main

import (
	"flag"
	"fmt"
	"log"

	"yandex-catalog-handler/internal/concator"
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

	consumerService := consumer.New(cfg, loaderService)

	s := server.NewServer(cfg, consumerService)

	//err := loaderService.Load()
	//
	//if err != nil {
	//	log.Printf("Error when try loaderService.Load, err: ", err)
	//}
	//return

	concatorService := concator.New(cfg)

	result, err := concatorService.Concate()

	if err != nil {
		log.Printf("Error when try concatorService.Concate, err: ", err)
		return
	}

	for _, elem := range result.Results {
		fmt.Println(elem)
	}

	return

	err = s.Run()

	if err != nil {
		log.Fatal("Error when try server.NewServer, err: %s", err)
	}
}
