package result

import (
	"encoding/json"
	"yandex-catalog-handler/internal/entity"
	"yandex-catalog-handler/pkg/storage"
)

type Repository struct {
	db *storage.Storage
}

func NewRepository(db *storage.Storage) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(result entity.Result) (err error) {
	jsonResults, err := json.Marshal(result.Results)

	if err != nil {
		return
	}

	_, err = r.db.Exec(
		"INSERT INTO operate_log(cause, results) VALUES ($1, $2)",
		result.Cause, jsonResults)

	return err
}
