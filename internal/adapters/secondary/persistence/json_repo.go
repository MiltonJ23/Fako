package persistence

import (
	"encoding/json"
	"os"

	"github.com/MiltonJ23/Fako/internal/core/domain"
)

type JSONStateRepository struct {
	Filepath string
}

func NewJSONStateRepository(path string) *JSONStateRepository {
	return &JSONStateRepository{Filepath: path}
}

func (j *JSONStateRepository) Load() (*domain.State, error) {
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
	// The first thing we are going to do is to marshall the state, here it is the marshalling operation
	data, marshallingError := json.MarshalIndent(state, "", "		")
	if marshallingError != nil {
		return marshallingError // Terrible but true to the most damn students of the hall
	}
	return os.WriteFile(j.Filepath, data, 0644)
}
