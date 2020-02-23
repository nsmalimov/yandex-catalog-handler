package result

import (
	"yandex-catalog-handler/internal/entity"
)

type Repo interface {
	Create(result entity.Result) error
	GetAll() ([]*entity.Result, error)
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{repo}
}

func (s *Service) Create(result entity.Result) (err error) {
	err = s.repo.Create(result)

	return
}

func (s *Service) GetAll() (results []*entity.Result, err error) {
	results, err = s.repo.GetAll()

	return
}
