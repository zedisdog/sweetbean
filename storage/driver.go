package storage

type Driver interface {
	Put(path string, data []byte) error
	Get(path string) ([]byte, error)
	Remove(path string) error
}

type DriverHasMime interface {
	MimeType(path string) string
}

type DriverHasPath interface {
	//Path 绝对路径
	Path(path string) string
}

type DriverHasBase64 interface {
	Base64(path string) (string, error)
}

type DriverCanGetSize interface {
	Size(path string) (int, error)
}

type DriverHasUrl interface {
	Url(path string) string
}
