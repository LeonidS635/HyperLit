package index

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"

	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

const indexFile = "index"
const indexFilePerm = 0644

type Index struct {
	path string
}

func NewIndex(path string) Index {
	return Index{path: filepath.Join(path, indexFile)}
}

func (i Index) Init() error {
	_, err := os.OpenFile(i.path, os.O_CREATE|os.O_EXCL, indexFilePerm)
	return err
}

func (i Index) Read(entriesInfoCh chan<- entry.Info) error {
	file, err := os.Open(i.path)
	if err != nil {
		return err
	}

	defer file.Close()
	defer close(entriesInfoCh)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		info, err := entry.ParseInfo(line)
		if err != nil {
			return err
		}
		entriesInfoCh <- info
	}
	return scanner.Err()
}

func (i Index) Save(data []entry.Info) error {
	file, err := os.Open(i.path)
	if err != nil {
		return err
	}
	defer file.Close()

	sort.Slice(
		data, func(i, j int) bool {
			return data[i].Path < data[j].Path
		},
	)

	writer := bufio.NewWriter(file)
	for _, info := range data {
		_, err := writer.Write(info.Dump())
		if err != nil {
			return err
		}
	}

	return nil
}
