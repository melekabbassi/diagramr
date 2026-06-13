package output

import (
	"io"
	"os"
)

func Write(path, content string, w io.Writer) error {
	if path == "" {
		_, err := io.WriteString(w, content)
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
