package drivers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/h2non/filetype"
	"github.com/zedisdog/sweetbean/storage"
	"github.com/zedisdog/sweetbean/tools"
)

func WithBaseUrl(baseUrl string) func(*LocalDriver) {
	return func(ld *LocalDriver) {
		ld.baseUrl = strings.TrimRight(baseUrl, "/")
	}
}

func NewLocal(path *tools.Path, options ...func(*LocalDriver)) *LocalDriver {
	driver := &LocalDriver{
		root: path,
		perm: 0755,
	}
	for _, set := range options {
		set(driver)
	}
	return driver
}

var _ storage.Driver = (*LocalDriver)(nil)
var _ storage.DriverHasMime = (*LocalDriver)(nil)
var _ storage.DriverHasPath = (*LocalDriver)(nil)
var _ storage.DriverHasBase64 = (*LocalDriver)(nil)
var _ storage.DriverCanGetSize = (*LocalDriver)(nil)
var _ storage.DriverHasUrl = (*LocalDriver)(nil)

type LocalDriver struct {
	root    *tools.Path
	perm    fs.FileMode
	baseUrl string
}

func (l LocalDriver) Put(path string, data []byte) (err error) {
	err = os.MkdirAll(l.root.Dir(path), l.perm)
	if err != nil {
		return
	}
	f, err := os.OpenFile(l.root.Concat(path), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, l.perm)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = f.Write(data)
	return
}

func (l LocalDriver) Get(path string) (data []byte, err error) {
	f, err := os.Open(l.root.Concat(path))
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	data, err = io.ReadAll(f)
	return
}

func (l LocalDriver) Remove(path string) (err error) {
	_, err = os.Stat(l.root.Concat(path))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return
	}

	err = os.Remove(l.root.Concat(path))
	return
}

func (l LocalDriver) Path(path string) string {
	return l.root.Concat(path)
}

func (l LocalDriver) MimeType(path string) string {
	fp, err := os.Open(l.root.Concat(path))
	if err != nil {
		return ""
	}
	defer func() {
		_ = fp.Close()
	}()
	b := make([]byte, 262)
	if _, err := fp.Read(b); err != nil {
		return ""
	}
	kind, _ := filetype.Match(b)
	if kind == filetype.Unknown {
		return ""
	}

	return kind.MIME.Value
}

func (l LocalDriver) Base64(path string) (s string, err error) {
	d, err := l.Get(path)
	if err != nil {
		return
	}
	s = base64.StdEncoding.EncodeToString(d)
	return
}

func (l LocalDriver) Size(path string) (size int, err error) {
	d, err := l.Get(path)
	if err != nil {
		return
	}
	size = len(d)
	return
}

func (l LocalDriver) Url(path string) string {
	return fmt.Sprintf("%s/%s", l.baseUrl, path)
}
