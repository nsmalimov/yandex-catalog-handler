package main

import (
	"flag"
	"fmt"
	"log"

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

	fmt.Println(cfg)
}
