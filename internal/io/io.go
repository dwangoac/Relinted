package io

import "os"

// ReadFile reads the entire contents of the file at path.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to the file at path.
func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
