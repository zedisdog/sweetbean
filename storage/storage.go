package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/zedisdog/sweetbean/errx"
)

func NewStorage(driver Driver) *Storage {
	return &Storage{
		driver: driver,
	}
}

type Storage struct {
	driver Driver
}

//PutFileQuick is similar than PutFile, but don't set filename.
func (s Storage) PutFileQuick(file *multipart.FileHeader, directory string) (path string, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}
	//path = directory/xxxx.jpg
	path = fmt.Sprintf(
		"%s%s%s%s",
		strings.Trim(directory, "\\/"),
		string(os.PathSeparator),
		id.String(),
		filepath.Ext(file.Filename),
	)
	err = s.PutFile(path, file)
	return
}

func (s Storage) Put(path string, data []byte) error {
	return s.driver.Put(path, data)
}

func (s Storage) Get(path string) ([]byte, error) {
	return s.driver.Get(path)
}

func (s Storage) Remove(path string) error {
	return s.driver.Remove(path)
}

func (s Storage) PutString(path string, data string) error {
	return s.driver.Put(path, []byte(data))
}

func (s Storage) GetString(path string) (data string, err error) {
	tmp, err := s.driver.Get(path)
	if err != nil {
		return
	}
	data = string(tmp)
	return
}

func (s Storage) PutFile(path string, file *multipart.FileHeader) (err error) {
	fp, err := file.Open()
	if err != nil {
		return
	}
	defer func() {
		_ = fp.Close()
	}()
	data, err := io.ReadAll(fp)
	if err != nil {
		return
	}
	return s.Put(path, data)
}

func (s Storage) MimeType(path string) (string, error) {
	if ss, ok := interface{}(s.driver).(DriverHasMime); ok {
		return ss.MimeType(path), nil
	}
	return "", errx.New("driver is not implement interface <DriverHasMime>", 1)
}

func (s Storage) Path(path string) (string, error) {
	if ss, ok := interface{}(s.driver).(DriverHasPath); ok {
		return ss.Path(path), nil
	}
	return "", errx.New("driver is not implement interface <DriverHasPath>", 1)
}

func (s Storage) Base64(path string) (string, error) {
	if ss, ok := interface{}(s.driver).(DriverHasBase64); ok {
		return ss.Base64(path)
	}
	return "", errx.New("driver is not implement interface <DriverHasBase64>", 1)
}

func (s Storage) Size(path string) (int, error) {
	if ss, ok := interface{}(s.driver).(DriverCanGetSize); ok {
		return ss.Size(path)
	}
	return 0, errx.New("driver is not implement interface <DriverCanGetSize>", 1)
}

func (s Storage) Url(path string) (string, error) {
	if ss, ok := interface{}(s.driver).(DriverHasUrl); ok {
		return ss.Url(path), nil
	}
	return "", errx.New("driver is not implement interface <DriverHasUrl>", 1)
}
