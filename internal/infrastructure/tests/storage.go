package tests

import (
	"github.com/EugeneNail/actum/internal/infrastructure/env"
	"os"
	"path/filepath"
	"testing"
)

type Storage struct {
	path string
	t    *testing.T
}

func NewStorage(t *testing.T) Storage {
	return Storage{
		filepath.Join(env.Get("APP_PATH"), "storage"),
		t,
	}
}

func (storage *Storage) AssertCount(directory string, count int) *Storage {
	directory = filepath.Join(storage.path, directory)
	files, err := os.ReadDir(directory)
	Check(err)

	if len(files) != count {
		storage.t.Errorf(
			"Expected directory %s to have %d files, got %d",
			directory, count, len(files),
		)
		storage.t.SkipNow()
	}

	return storage
}

func (storage *Storage) AssertHas(directory string, filename string) *Storage {
	directory = filepath.Join(storage.path, directory)

	if !storage.hasFile(directory, filename) {
		storage.t.Errorf("Expected directory %s to have file %s", directory, filename)
	}

	return storage
}

func (storage *Storage) AssertLacks(directory string, filename string) *Storage {
	directory = filepath.Join(storage.path, directory)

	if storage.hasFile(directory, filename) {
		storage.t.Errorf("Expected directory %s not to have file %s", directory, filename)
	}

	return storage
}

func (storage *Storage) hasFile(directory string, filename string) bool {
	files, err := os.ReadDir(directory)
	Check(err)

	for _, file := range files {
		if file.Name() == filename {
			return true
		}
	}

	return false
}
