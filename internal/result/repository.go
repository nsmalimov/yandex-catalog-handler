package result

import (
	"yandex-catalog-handler/internal/entity"
	"yandex-catalog-handler/pkg/storage"
)

type Repository struct {
	db *storage.Storage
}

func NewRepository(db *storage.Storage) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(result entity.Result) error {
	_, err := r.db.Exec(
		"INSERT INTO operate_log(cause, results) VALUES (?, ?)",
		result.Cause, result.Results)

	return err
}
