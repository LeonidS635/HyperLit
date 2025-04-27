package roothash

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
)

const rootHashFileName = "root"

type RootHash struct {
	path string
}

func NewRoot(path string) RootHash {
	return RootHash{path: filepath.Join(path, rootHashFileName)}
}

func (r RootHash) Save(hash []byte) error {
	return os.WriteFile(r.path, hash, 0644)
}

func (r RootHash) Get() (string, error) {
	hash, err := os.ReadFile(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return hasher.ConvertToHex(hash), nil
}
