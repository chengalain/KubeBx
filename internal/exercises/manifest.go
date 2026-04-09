package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

type Exercise struct {
	Manifest
	Path string
}

// GetExercisesDir returns the path to the exercises directory
func GetExercisesDir() (string, error) {
	// Try current directory first
	if _, err := os.Stat("exercises"); err == nil {
		return "exercises", nil
	}

	// Try relative to executable
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	exercisesPath := filepath.Join(exPath, "exercises")

	if _, err := os.Stat(exercisesPath); err == nil {
		return exercisesPath, nil
	}

	return "", fmt.Errorf("exercises directory not found")
}

// LoadAll loads all exercises from the exercises directory
func LoadAll() ([]Exercise, error) {
	exercisesDir, err := GetExercisesDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(exercisesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read exercises directory: %w", err)
	}

	var exercises []Exercise

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(exercisesDir, entry.Name(), "manifest.yaml")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			// Skip directories without manifest
			continue
		}

		var manifest Manifest
		if err := yaml.Unmarshal(data, &manifest); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: invalid manifest in %s: %v\n", entry.Name(), err)
			continue
		}

		exercises = append(exercises, Exercise{
			Manifest: manifest,
			Path:     filepath.Join(exercisesDir, entry.Name()),
		})
	}

	// Sort by ID
	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].ID < exercises[j].ID
	})

	return exercises, nil
}

// FindByID finds an exercise by its ID
func FindByID(id string) (*Exercise, error) {
	exercises, err := LoadAll()
	if err != nil {
		return nil, err
	}

	for _, ex := range exercises {
		if ex.ID == id {
			return &ex, nil
		}
	}

	return nil, fmt.Errorf("exercise %s not found", id)
}
