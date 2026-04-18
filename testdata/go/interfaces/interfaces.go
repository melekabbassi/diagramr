package interfaces

type Reader interface {
	Read() string
}

type FileReader struct{}

func (f *FileReader) Read() string {
	return "ok"
}
