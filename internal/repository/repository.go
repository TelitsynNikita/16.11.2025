package repository

type URLStorage interface {
	GetUrlByID(id uint) (int, error)
}

type Repository struct {
	URLStorage
}

func NewRepository() *Repository {
	return &Repository{
		URLStorage: NewPersistentURLStorage(),
	}
}
