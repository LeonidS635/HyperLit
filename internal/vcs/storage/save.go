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

func (s ObjectsStorage) SaveEntryTmp(entry entry.Interface) error {
	return s.saveDataTmp(hasher.ConvertToHex(entry.GetHash()), entry.GetData())
}

func (s ObjectsStorage) saveDataTmp(hash string, data []byte) error {
	dirName := filepath.Join(s.tmpDir, hash[:2])
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

func (s ObjectsStorage) Dump() error {
	dirs, err := os.ReadDir(s.tmpDir)
	if err != nil {
		return err
	}

	for _, d := range dirs {
		files, err := os.ReadDir(filepath.Join(s.tmpDir, d.Name()))
		if err != nil {
			return err
		}

		dirTmpPath := filepath.Join(s.tmpDir, d.Name())
		dirDestPath := filepath.Join(s.workingDir, d.Name())

		if err = os.Mkdir(dirDestPath, dirPermissions); err != nil && !os.IsExist(err) {
			return err
		}

		for _, file := range files {
			fileTmpPath := filepath.Join(dirTmpPath, file.Name())
			fileDestPath := filepath.Join(dirDestPath, file.Name())

			if err = os.Rename(fileTmpPath, fileDestPath); err != nil {
				return err
			}
		}
	}

	//return os.RemoveAll(s.tmpDir)
	return nil
}
