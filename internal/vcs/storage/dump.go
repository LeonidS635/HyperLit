package storage

import (
	"fmt"
	"os"
)

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
