package persistence

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type repository struct {
	sync.RWMutex
	data map[int]string
}

func (repo *repository) Insert(encodedHash string) (key int) {
	repo.Lock()
	defer repo.Unlock()
	position := len(repo.data) + 1
	repo.data[position] = encodedHash
	return position
}

func (repo *repository) Get(position int) (value string, err error) {
	repo.RLock()
	defer repo.RUnlock()

	var ok bool

	if value, ok = repo.data[position]; ok != true {
		err = fmt.Errorf("cannot access data at position %d", position)
	}

	return value, err

}

func (repo *repository) Update(key int, value string) {
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
			data: make(map[int]string),
		}

		atomic.StoreInt32(&initialized, 1)
	}

	return instance

}
