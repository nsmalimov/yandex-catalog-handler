package consumer

import (
	"log"
	"time"

	"yandex-catalog-handler/internal/loader"
	"yandex-catalog-handler/pkg/config"
)

type Consumer struct {
	tt     *time.Timer
	delta  int
	cfg    config.Config
	loader *loader.Loader
}

func New(cfg config.Config, loader *loader.Loader) *Consumer {
	return &Consumer{
		cfg:    cfg,
		loader: loader,
	}
}

func (c *Consumer) Run() {
	c.tt = time.NewTimer(time.Duration(c.delta) * time.Second)

	for {
		select {
		case _, ok := <-c.tt.C:
			if !ok {
				return
			}
			log.Printf("Start load data")
			c.tt.Reset(time.Duration(c.delta) * time.Second)
		}
	}
}
