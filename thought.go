package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type thought struct {
	Text       string    `json:"text"`
	FlowerType int       `json:"flower_type"`
	ColorIndex int       `json:"color_index"`
	CreatedAt  time.Time `json:"created_at"`
}

func thoughtsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "thought-garden", "thoughts.json"), nil
}

func loadThoughts() ([]thought, error) {
	path, err := thoughtsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []thought{}, nil
	}
	if err != nil {
		return nil, err
	}
	var thoughts []thought
	if err := json.Unmarshal(data, &thoughts); err != nil {
		return nil, err
	}
	return thoughts, nil
}

func saveThoughts(thoughts []thought) error {
	path, err := thoughtsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(thoughts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}
