package persistence

import (
	"fmt"
	"sync"
)

//Repository is a simple, thread-safe, keyed data store allowing for efficient inserts and lookups
type Repository struct {
	sync.RWMutex
	data map[int64]string
}

//Insert adds a string value to the data store and returns a key for later retrieval
func (repo *Repository) Insert(encodedHash string) (key int64) {
	repo.Lock()
	defer repo.Unlock()

	position := int64(len(repo.data) + 1)
	repo.data[position] = encodedHash
	return position
}


//Get retrieves a string value for a given key
func (repo *Repository) Get(key int64) (value string, err error) {
	repo.RLock()
	defer repo.RUnlock()

	var ok bool

	if value, ok = repo.data[key]; ok != true {
		err = fmt.Errorf("cannot access data for key %d", key)
	}

	return value, err
}

//Update overwrites the value at a given key
func (repo *Repository) Update(key int64, value string) {
	repo.Lock()
	defer repo.Unlock()

	repo.data[key] = value
}

//NewRepository instantiates and returns a new instance of Repository
func NewRepository() *Repository {
	return &Repository{
		data: make(map[int64]string),
	}
}
