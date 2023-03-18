package dir

import (
	"os"
	"path/filepath"
)

func CreateDir(dir string) (*string, error) {
	result := filepath.Join("kup", dir)
	if err := os.MkdirAll(result, os.ModePerm); err != nil {
		return nil, err
	}
	return &result, nil
}
