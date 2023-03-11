package zip

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

func ZipWriter(baseFolder, destFileName string) {

	// Get a Buffer to Write To
	outFile, err := os.Create(destFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "")

	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		filePath := filepath.Join(basePath, file.Name())
		if !file.IsDir() {
			dat, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + filePath)

			addFiles(w, filePath, filepath.Join(baseInZip, file.Name()))
		}
	}
}
