package persistence

import (
	"fmt"
	"math/rand"
	"time"
)

func newRepo() *repository {
	return &repository{data: map[uint64]string{}}
}

func ExampleRepository_InsertInEmptyRepositoryReturnsOne() {
	repo := newRepo()
	key := repo.Insert("Hello World")
	fmt.Println(key)
	// Output: 1
}

func ExampleRepository_InsertMultipleItemsReturnsSequentialKeys() {
	repo := newRepo()
	count := 5
	keys := make([]uint64, count)

	for i := 0; i < count; i++ {
		keys[i] = repo.Insert("test")
	}

	fmt.Println(keys)
	// Output: [1 2 3 4 5]

}

func ExampleRepository_ConcurrentInsertsGuaranteedSequentialKeys() {
	repo := newRepo()
	count := 5
	keys := make([]uint64, count)
	results := make(chan uint64)

	for i := 0; i < count; i++ {
		go func() {
			randomDuration := time.Duration(rand.Intn(10)) * time.Microsecond
			time.Sleep(randomDuration)
			results <- repo.Insert("test")
		}()
	}

	for i := 0; i < count; i++ {
		keys[i] = <-results
	}

	fmt.Println(keys)
	// Output: [1 2 3 4 5]

}

func ExampleRepository_GetAtInvalidPositionReturnsError() {
	repo := newRepo()

	_, err := repo.Get(1)
	fmt.Println(err)
	// Output: cannot access data at position 1
}

func ExampleRepository_GetAtValidPositionReturnsStoredValue() {
	repo := newRepo()
	key := repo.Insert("This is a test")

	value, _ := repo.Get(key)
	fmt.Println(value)
	// Output: This is a test
}

func ExampleRepository_GetAtValidPositionReturnsNoError() {
	repo := newRepo()
	key := repo.Insert("This is a test")

	_, err := repo.Get(key)
	fmt.Println(err)
	// Output: <nil>
}

func ExampleRepository_Update() {
	repo := newRepo()
	key := repo.Insert("Hello World")
	repo.Update(key, "Goodbye World")

	value, _ := repo.Get(key)
	fmt.Println(value)
	// Output: Goodbye World
}
