package persistence

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

type JSONStateRepository struct {
	Filepath string
	mu       sync.RWMutex
}

func NewJSONStateRepository(path string) *JSONStateRepository {
	return &JSONStateRepository{Filepath: path}
}

func (j *JSONStateRepository) Load() (*domain.State, error) {
	// let's first of all acquire the RLock , so that multiple goroutines can go ahead and load at the same time
	j.mu.RLock()
	defer j.mu.RUnlock()

	// first of all , let's check if the file do exist
	_, FileStatError := os.Stat(j.Filepath)
	if os.IsNotExist(FileStatError) {
		return domain.NewState(), nil
	}

	// we reached here, meaning the file do exist
	data, ReadingFileError := os.ReadFile(j.Filepath)
	if ReadingFileError != nil {
		return nil, ReadingFileError
	}
	// we made sure we could get the content of the file, if we reach here it is okay
	// We are then going to decode the json
	var state domain.State
	marshallingError := json.Unmarshal(data, &state) // This is the error for Unmarshalling operation no mtter how it is called
	if marshallingError != nil {
		return nil, marshallingError
	}
	return &state, nil
}

func (j *JSONStateRepository) Save(state *domain.State) error {
	// here we are going with the Write Lock which is clearly the most sensible operation
	j.mu.Lock()         // we are preventing corruption
	defer j.mu.Unlock() // automatically free the tex at the end of the execution

	// The first thing we are going to do is to marshall the state, here it is the marshalling operation
	data, marshallingError := json.MarshalIndent(state, "", "		")
	if marshallingError != nil {
		return marshallingError // Terrible but true to the most damn students of the hall
	}
	return os.WriteFile(j.Filepath, data, 0644)
}
