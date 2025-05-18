package storage

import (
	"fmt"
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

// TODO: do something with ctx (deal with cancel)

func (s *ObjectsStorage) Dump() error {
	var err error
	defer func() {
		if err != nil {
			for hash := range s.renames {
				if restoreErr := s.restoreOldData(hash); restoreErr != nil {
					err = fmt.Errorf("[FATAL] failed to restore %s: %w", hash, restoreErr)
					return
				}
			}
		}
	}()

	for hash := range s.renames {
		if err = s.saveOldData(hash); err != nil {
			return err
		}
	}

	if err = os.RemoveAll(s.workingDir); err != nil {
		return err
	}
	if err = os.Rename(s.tmpDir, s.workingDir); err != nil {
		return err
	}
	return nil
}

func (s *ObjectsStorage) Clear() error {
	return os.RemoveAll(s.tmpDir)
}
