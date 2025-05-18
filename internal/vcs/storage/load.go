package storage

import (
	"fmt"
	"os"

	"github.com/LeonidS635/HyperLit/internal/vcs/objects/blob"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (s *ObjectsStorage) LoadEntry(hash string) (entry.Interface, error) {
	_, filePath := getDirAndFilePathByHash(s.workingDir, hash)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	type_, _, err := format.ParseHeaderFromData(data)
	if err != nil {
		return nil, err
	}

	switch type_ {
	case format.TreeType:
		return tree.NewTree(data), nil
	case format.CodeType, format.DocsType:
		return blob.NewBlob(type_, data), nil
	default:
		return nil, fmt.Errorf("error loading entry: unknown type: %v", type_)
	}
}
