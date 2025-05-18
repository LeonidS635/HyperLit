package globs

import (
	"bufio"
	"os"
	"path/filepath"
)

const ignoreFileName = ".hlignore"

func GetGlobsToIgnore(projectPath string) []string {
	ignoreFile, err := os.Open(filepath.Join(projectPath, ignoreFileName))
	if err != nil && !os.IsNotExist(err) {
		return nil
	}
	defer ignoreFile.Close()

	var globs []string

	scanner := bufio.NewScanner(ignoreFile)
	for scanner.Scan() {
		line := scanner.Text()
		globs = append(globs, line)
	}
	return globs
}
