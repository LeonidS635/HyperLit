package storage

import "path/filepath"

func getDirAndFilePathByHash(rootPath, hash string) (string, string) {
	dirName := filepath.Join(rootPath, hash[:2])
	fileName := filepath.Join(dirName, hash[2:])
	return dirName, fileName
}
