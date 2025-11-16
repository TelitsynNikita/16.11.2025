package service

import (
	"workmate/internal/model"
	"workmate/internal/repository"
)

type PersistentURLService interface {
	GetUrlByID(ids []int) ([]model.CheckLinksStatusByUrlResponse, error)
	CheckLinksStatusByUrl(urls []string) (model.CheckLinksStatusByUrlResponse, error)
}

type Service struct {
	PersistentURLService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		PersistentURLService: NewURLService(repo.URLStorageRepository),
	}
}
