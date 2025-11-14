package service

import "workmate/internal/repository"

type URLService interface {
	GetUrlByID(id uint) (int, error)
}

type Service struct {
	URLService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		URLService: NewURLService(repo.URLStorage),
	}
}
