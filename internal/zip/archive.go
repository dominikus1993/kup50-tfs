package zip

import (
	"archive/zip"
	"log"
	"os"
)

func ArchiveAllFiles() {
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	archive, err := os.Create("archive.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
}
