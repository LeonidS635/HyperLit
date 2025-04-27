package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

func (s ObjectsStorage) SaveEntry(entry entry.Interface) error {
	return s.saveData(hasher.ConvertToHex(entry.GetHash()), entry.GetData())
}

func (s ObjectsStorage) saveData(hash string, data []byte) error {
	dirName := filepath.Join(s.workingDir, hash[:2])
	fileName := filepath.Join(dirName, hash[2:])

	err := os.Mkdir(dirName, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}

	err = os.WriteFile(fileName, data, filePermissions)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
