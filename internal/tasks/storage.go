// storage.go
package tasks

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const configDirName = ".config/bcli"
const tasksFileName = "tasks.json"

func getStoragePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDirName, tasksFileName), nil
}

func ensureStorageDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(filepath.Join(home, configDirName), os.ModePerm)
}

func SaveTasks(tasks []Task) error {
	if err := ensureStorageDir(); err != nil {
		return err
	}

	path, err := getStoragePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadTasks() ([]Task, error) {
	path, err := getStoragePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
