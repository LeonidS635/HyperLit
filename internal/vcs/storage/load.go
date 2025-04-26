package storage

import (
	"os"

	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
)

func (s ObjectsStorage) LoadEntry(hash string) (entry.Entry, error) {
	data, err := os.ReadFile(s.getFilePathByHash(hash))
	if err != nil {
		return entry.Entry{}, err
	}

	type_, size, err := format.ParseHeaderFromData(data)
	return entry.Entry{
		Type: type_,
		Size: size,
		Data: data[format.HeaderSize:],
	}, err
}
