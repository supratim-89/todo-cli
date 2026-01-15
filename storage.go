package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gofrs/flock"
)

const fileName = "todos.json"

var fileLock = flock.New(fileName + ".lock")

func loadTodos() ([]Todo, error) {
	// Acquire shared lock (read)
	if err := fileLock.RLock(); err != nil {
		return nil, err
	}
	defer fileLock.Unlock()

	// Auto-create file if missing
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := os.WriteFile(fileName, []byte("[]"), 0644)
		if err != nil {
			return nil, err
		}
		return []Todo{}, nil
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("todos.json is empty; expected []")
	}

	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func saveTodos(todos []Todo) error {
	// Acquire exclusive lock (write)
	if err := fileLock.Lock(); err != nil {
		return err
	}
	defer fileLock.Unlock()

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	// Atomic write
	tmp := fileName + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, fileName)
}
