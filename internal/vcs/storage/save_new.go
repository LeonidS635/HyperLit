package storage

import (
	"os"

	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

func (s *ObjectsStorage) SaveNewEntry(entry entry.Interface) error {
	return s.saveNewData(hasher.ConvertToHex(entry.GetHash()), entry.GetData())
}

func (s *ObjectsStorage) saveNewData(hash string, data []byte) error {
	dirPath, filePath := getDirAndFilePathByHash(s.tmpDir, hash)

	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Mkdir(dirPath, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return os.WriteFile(filePath, data, filePermissions)
}
