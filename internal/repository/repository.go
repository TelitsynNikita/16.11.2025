package repository

import "workmate/internal/model"

type URLStorageRepository interface {
	GetUrlByIDs(ids []int) ([]model.PersistentStorageData, error)
	GetLinksByUrl(urls []string) (int, []string, error)
	WriteDataToFileAndLocalStorage() error
	ReadFileToLocalStorage() error
	InitPersistentStorage() error
}

type Repository struct {
	URLStorageRepository
}

func NewRepository() *Repository {
	return &Repository{
		URLStorageRepository: NewPersistentURLStorage(),
	}
}
