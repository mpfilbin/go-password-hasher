package persistence

import (
	"fmt"
	"sync"
)

type Repository struct {
	sync.RWMutex
	data map[int64]string
}

func (repo *Repository) Insert(encodedHash string) (key int64) {
	repo.Lock()
	defer repo.Unlock()
	position := int64(len(repo.data) + 1)
	repo.data[position] = encodedHash
	return position
}

func (repo *Repository) Get(position int64) (value string, err error) {
	repo.RLock()
	defer repo.RUnlock()

	var ok bool

	if value, ok = repo.data[position]; ok != true {
		err = fmt.Errorf("cannot access data at position %d", position)
	}

	return value, err

}

func (repo *Repository) Update(key int64, value string) {
	repo.Lock()
	defer repo.Unlock()

	repo.data[key] = value
}

func NewRepository() *Repository {
	return &Repository{
		data: make(map[int64]string),
	}
}