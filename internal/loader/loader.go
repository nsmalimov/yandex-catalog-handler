package loader

import "yandex-catalog-handler/pkg/config"

type Loader struct {
	cfg config.Config
}

func New(cfg config.Config) *Loader {
	return &Loader{
		cfg: cfg,
	}
}
