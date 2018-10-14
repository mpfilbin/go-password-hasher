package persistence

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type repository struct {
	sync.RWMutex
	data map[uint64]string
}

func (repo *repository) Insert(encodedHash string) (key uint64) {
	repo.Lock()
	defer repo.Unlock()
	position := uint64(len(repo.data) + 1)
	repo.data[position] = encodedHash
	return position
}

func (repo *repository) Get(position uint64) (value string, err error) {
	repo.RLock()
	defer repo.RUnlock()

	var ok bool

	if value, ok = repo.data[position]; ok != true {
		err = fmt.Errorf("cannot access data at position %d", position)
	}

	return value, err

}

func (repo *repository) Update(key uint64, value string) {
	repo.Lock()
	defer repo.Unlock()

	repo.data[key] = value
}

var instance *repository
var mu sync.Mutex
var initialized int32

func GetInstance() *repository {

	if atomic.LoadInt32(&initialized) == 1 {
		return instance
	}

	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = &repository{
			data: make(map[uint64]string),
		}

		atomic.StoreInt32(&initialized, 1)
	}

	return instance

}
