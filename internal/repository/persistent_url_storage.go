package repository

type PersistentURLStorage struct{}

func NewPersistentURLStorage() *PersistentURLStorage {
	return &PersistentURLStorage{}
}

func (p *PersistentURLStorage) GetUrlByID(id uint) (int, error) {
	return 0, nil
}
