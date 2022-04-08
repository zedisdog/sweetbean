package drivers

import (
	"encoding/base64"
	"errors"
	"github.com/h2non/filetype"
	"github.com/zedisdog/sweetbean/tools"
	"io"
	"io/fs"
	"os"
)

func NewLocal(path *tools.Path) *LocalDriver {
	return &LocalDriver{
		path: path,
		perm: 0755,
	}
}

type LocalDriver struct {
	path *tools.Path
	perm fs.FileMode
}

func (l LocalDriver) Put(path string, data []byte) (err error) {
	err = os.MkdirAll(l.path.Dir(path), l.perm)
	if err != nil {
		return
	}
	f, err := os.OpenFile(l.path.Concat(path), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, l.perm)
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
	f, err := os.Open(l.path.Concat(path))
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
	_, err = os.Stat(l.path.Concat(path))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return
	}

	err = os.Remove(l.path.Concat(path))
	return
}

func (l LocalDriver) Path(path string) string {
	return l.path.Concat(path)
}

func (l LocalDriver) MimeType(path string) string {
	fp, err := os.Open(l.path.Concat(path))
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
