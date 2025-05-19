package globs

import (
	"bufio"
	"os"
	"path/filepath"
)

const ignoreFileName = ".hlignore"

func GetGlobsToIgnore(projectPath string) []string {
	ignoreFilePath := filepath.Join(projectPath, ignoreFileName)
	ignoreFile, err := os.Open(ignoreFilePath)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}
	defer ignoreFile.Close()

	var globs []string
	globs = append(globs, ignoreFilePath)

	scanner := bufio.NewScanner(ignoreFile)
	for scanner.Scan() {
		line := scanner.Text()
		globs = append(globs, filepath.Join(projectPath, line))
	}
	return globs
}
