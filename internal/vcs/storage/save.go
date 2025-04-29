package storage

import (
	"fmt"
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
	fmt.Println("Saving new", hash)
	dirPath, filePath := GetDirAndFilePathByHash(s.tmpDir, hash)

	err := os.Mkdir(dirPath, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return os.WriteFile(filePath, data, filePermissions)
}

func (s ObjectsStorage) saveOldData(hash string) error {
	fmt.Println("Saving old", hash)
	dirPath, filePath := GetDirAndFilePathByHash(s.tmpDir, hash)
	_, origFilePath := GetDirAndFilePathByHash(s.workingDir, hash)

	err := os.Mkdir(dirPath, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return os.Rename(origFilePath, filePath)
}

func (s ObjectsStorage) Dump() error {
	//dirs, err := os.ReadDir(s.tmpDir)
	//if err != nil {
	//	return err
	//}
	//
	//for _, dir := range dirs {
	//	files, err := os.ReadDir(filepath.Join(s.tmpDir, dir.Name()))
	//	if err != nil {
	//		return err
	//	}
	//
	//	if err := os.Mkdir(filepath.Join(s.workingDir, dir.Name()), dirPermissions); err != nil && !os.IsExist(err) {
	//		return err
	//	}
	//
	//	for _, file := range files {
	//		if err := os.Rename(
	//			filepath.Join(s.tmpDir, dir.Name(), file.Name()), filepath.Join(s.workingDir, dir.Name(), file.Name()),
	//		); err != nil {
	//			fmt.Println("$$$$$$$$$$$$$$$$$$$")
	//			fmt.Println(err)
	//			fmt.Println("$$$$$$$$$$$$$$$$$$$")
	//			return err
	//		}
	//	}
	//}
	//
	//return os.RemoveAll(s.tmpDir)
	//os.Rename(s.workingDir, fmt.Sprintf("%s_backup", s.workingDir))
	if err := os.RemoveAll(s.workingDir); err != nil {
		return err
	}
	return os.Rename(s.tmpDir, s.workingDir)
	//return nil
	//dirs, err := os.ReadDir(s.tmpDir)
	//if err != nil {
	//	return err
	//}
	//
	//for _, d := range dirs {
	//	files, err := os.ReadDir(filepath.Join(s.tmpDir, d.Name()))
	//	if err != nil {
	//		return err
	//	}
	//
	//	dirTmpPath := filepath.Join(s.tmpDir, d.Name())
	//	dirDestPath := filepath.Join(s.workingDir, d.Name())
	//
	//	if err = os.Mkdir(dirDestPath, dirPermissions); err != nil && !os.IsExist(err) {
	//		return err
	//	}
	//
	//	for _, file := range files {
	//		fileTmpPath := filepath.Join(dirTmpPath, file.Name())
	//		fileDestPath := filepath.Join(dirDestPath, file.Name())
	//
	//		if err = os.Rename(fileTmpPath, fileDestPath); err != nil {
	//			return err
	//		}
	//	}
	//}
	//
	//return os.RemoveAll(s.tmpDir)
	//return nil
}

func (s ObjectsStorage) Clear() error {
	return os.RemoveAll(s.tmpDir)
}
