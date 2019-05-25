package helper

import (
	"path/filepath"
	"os"
)

func GetCurrentDir() string{
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		//log.Fatal(err)
		panic(err)
	}
	return dir
}
