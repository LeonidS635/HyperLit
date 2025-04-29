package storage

import (
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

func GetDirAndFilePathByHash(rootPath, hash string) (string, string) {
	dirName := filepath.Join(rootPath, hash[:2])
	fileName := filepath.Join(dirName, hash[2:])
	return dirName, fileName
}

func (s ObjectsStorage) SaveNewEntry(entry entry.Interface) error {
	return s.saveNewData(hasher.ConvertToHex(entry.GetHash()), entry.GetData())
}

func (s ObjectsStorage) SaveOldEntry(entry entry.Interface) error {
	return s.saveOldData(hasher.ConvertToHex(entry.GetHash()))
}

func (s ObjectsStorage) saveNewData(hash string, data []byte) error {
	dirPath, filePath := GetDirAndFilePathByHash(s.tmpDir, hash)

	err := os.Mkdir(dirPath, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return os.WriteFile(filePath, data, filePermissions)
}

func (s ObjectsStorage) saveOldData(hash string) error {
	dirPath, filePath := GetDirAndFilePathByHash(s.tmpDir, hash)
	_, origFilePath := GetDirAndFilePathByHash(s.workingDir, hash)

	err := os.Mkdir(dirPath, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return os.Rename(origFilePath, filePath)
}

func (s ObjectsStorage) Dump() error {
	if err := os.RemoveAll(s.workingDir); err != nil {
		return err
	}
	return os.Rename(s.tmpDir, s.workingDir)
}

func (s ObjectsStorage) Clear() error {
	return os.RemoveAll(s.tmpDir)
}
