package utils

import (
	"os"
	"path/filepath"

	"github.com/zedisdog/sweetbean/errx"
)

func CreateFile(path string, content []byte) (err error) {
	_, err = os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		err = errx.Wrap(err, "get file stat error")
		return
	}

	if !os.IsNotExist(err) {
		err = errx.New("file is exists")
		return
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		err = errx.Wrap(err, "mkdir error")
		return
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(errx.Wrap(err, "open file error"))
	}
	defer file.Close()
	_, err = file.Write(content)
	return
}
