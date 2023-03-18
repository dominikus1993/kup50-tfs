package dir

import (
	"errors"
	"io"
	"os"
)

var errNilContent = errors.New("can't save nil contnet")

func Save(filename string, content io.ReadCloser) error {
	defer content.Close()
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	// handle err
	_, err = io.Copy(outFile, content)
	// handle err
	if err != nil {
		return err
	}
	return nil
}

func SaveString(filename string, content *string) error {
	if content == nil {
		return errNilContent
	}
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	// handle err
	_, err = io.WriteString(outFile, *content)
	// handle err
	if err != nil {
		return err
	}
	return nil
}
