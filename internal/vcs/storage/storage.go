package storage

import (
	"os"
	"path/filepath"
)

type ObjectsStorage struct {
	workingDir string
}

func NewObjectsStorage(path string) ObjectsStorage {
	return ObjectsStorage{workingDir: filepath.Join(path, objectsDirName)}
}

func (s ObjectsStorage) Init() error {
	err := os.Mkdir(s.workingDir, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

//func (s *ObjectsStorage) Delete(hash string) error {
//	dirName, fileName := s.getDirAndFilePathsByHash(hash)
//	return os.Remove(filepath.Join(dirName, fileName))
//}

func (s ObjectsStorage) getFilePathByHash(hash string) string {
	dirName := filepath.Join(s.workingDir, hash[:2])
	fileName := filepath.Join(dirName, hash[2:])
	return fileName
}
