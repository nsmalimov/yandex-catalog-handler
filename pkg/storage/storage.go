package storage

import (
	"fmt"

	"yandex-catalog-handler/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	*sqlx.DB
}

func New(cfg config.Config) (*Storage, error) {
	db, err := sqlx.Connect("postgres", createDSN(cfg))
	if err != nil {
		return nil, err
	}
	s := &Storage{}
	s.DB = db

	return s, err
}

func createDSN(cfg config.Config) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", cfg.Db.UserName, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.DatabaseName)
}
