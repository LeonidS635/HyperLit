package storage

import (
	"os"
)

func (s *ObjectsStorage) SaveOldEntry(hash string) error {
	s.renamesLock.Lock()
	defer s.renamesLock.Unlock()

	s.renames[hash] = struct{}{}

	return nil
}

func (s *ObjectsStorage) saveOldData(hash string) error {
	dirPath, filePath := getDirAndFilePathByHash(s.tmpDir, hash)
	_, origFilePath := getDirAndFilePathByHash(s.workingDir, hash)

	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Mkdir(dirPath, dirPermissions)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return os.Rename(origFilePath, filePath)
}

func (s *ObjectsStorage) restoreOldData(hash string) error {
	_, filePath := getDirAndFilePathByHash(s.tmpDir, hash)
	_, origFilePath := getDirAndFilePathByHash(s.workingDir, hash)

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.Rename(filePath, origFilePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
