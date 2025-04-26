package dirtraverser

import (
	"os"
	"time"
)

type FileInfo struct {
	Path    string
	Mode    os.FileMode
	ModTime time.Time
}
