package service

import "workmate/internal/repository"

type URL struct {
	URLStorage repository.URLStorage
}

func NewURLService(repo repository.URLStorage) *URL {
	return &URL{
		URLStorage: repo,
	}
}

func (s *URL) GetUrlByID(id uint) (int, error) {
	return 0, nil
}
