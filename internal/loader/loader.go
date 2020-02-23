package loader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"yandex-catalog-handler/pkg/config"
)

type Loader struct {
	cfg config.Config
}

func New(cfg config.Config) *Loader {
	return &Loader{
		cfg: cfg,
	}
}

func downloadFile(filepath string, url string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
		return
	}()

	out, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer func() {
		err = out.Close()
		return
	}()

	_, err = io.Copy(out, resp.Body)
	return
}

func (l *Loader) Load() (err error) {
	for _, fileName := range l.cfg.FileNames {
		filePath := fmt.Sprintf("%s/%s.xml", l.cfg.DataPath, fileName)
		url := fmt.Sprintf("%s%s", l.cfg.SourceUrl, fileName)

		log.Printf("Start downloading, %s", fileName)

		err = downloadFile(filePath, url)

		if err != nil {
			return
		}

		log.Printf("End downloading, %s", fileName)
	}

	return
}
