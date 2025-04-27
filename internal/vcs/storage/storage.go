package storage

import (
	"os"
	"path/filepath"
)

type ObjectsStorage struct {
	tmpDir     string
	workingDir string
}

func NewObjectsStorage(path string) ObjectsStorage {
	tmpPath, _ := os.MkdirTemp(path, objectsDirName)
	return ObjectsStorage{
		tmpDir:     tmpPath,
		workingDir: filepath.Join(path, objectsDirName),
	}
}

func (s ObjectsStorage) Init() error {
	err := os.Mkdir(s.workingDir, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (s ObjectsStorage) Delete(hash string) error {
	filePath := s.getFilePathByHash(hash)
	return os.Remove(filePath)
}

func (s ObjectsStorage) getFilePathByHash(hash string) string {
	dirName := filepath.Join(s.workingDir, hash[:2])
	fileName := filepath.Join(dirName, hash[2:])
	return fileName
}
