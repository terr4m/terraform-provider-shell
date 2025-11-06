package shell

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// GetMetaFilePath returns the path to the meta file.
func GetMetaFilePath() (string, error) {
	return getTempFile("tf-script-meta-*.json")
}

// GetOutFilePath returns the path to the output file.
func GetOutFilePath() (string, error) {
	return getTempFile("tf-script-output-*.json")
}

// GetErrorFilePath returns the path to the output file.
func GetErrorFilePath() (string, error) {
	return getTempFile("tf-script-error-*")
}

// ReadJSON reads a file as JSON and returns the contents.
func ReadJSON(p string) (any, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if !json.Valid(b) {
		return nil, fmt.Errorf("file is not valid JSON")
	}

	var r any
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// getTempFile creates a temporary file and returns the path.
func getTempFile(pattern string) (string, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
