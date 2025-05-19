package storage

import (
	"os"
	"path/filepath"
	"sync"
)

type ObjectsStorage struct {
	hlPath      string
	projectPath string

	tmpDir     string
	workingDir string

	mu sync.Mutex

	renamesLock sync.Mutex
	renames     map[string]struct{} // Old entries hashes
}

func NewObjectsStorage(projectPath, hlPath string) *ObjectsStorage {
	return &ObjectsStorage{
		hlPath:      hlPath,
		projectPath: projectPath,
		workingDir:  filepath.Join(hlPath, objectsDirName),
		renames:     make(map[string]struct{}),
	}
}

func (s *ObjectsStorage) Init() error {
	tmpPath, err := os.MkdirTemp(s.hlPath, objectsDirName)
	if err != nil {
		return err
	}
	s.tmpDir = tmpPath

	if err = os.Mkdir(s.workingDir, dirPermissions); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (s *ObjectsStorage) Clear() error {
	return os.RemoveAll(s.tmpDir)
}
