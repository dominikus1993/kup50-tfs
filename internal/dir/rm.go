package dir

import "os"

func Rm(dir string) error {
	return os.RemoveAll(dir)
}
