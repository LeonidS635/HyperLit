package tree

import (
	"bytes"
	"errors"

	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
)

type Tree struct {
	name    string
	content []byte
}

func Prepare(name string) (*Tree, error) {
	header, err := format.FormHeader(format.TreeType)
	if err != nil {
		return nil, err
	}

	return &Tree{
		name:    name,
		content: header,
	}, nil
}

func FromEntry(entry entry.Entry) (*Tree, error) {
	t, err := Prepare(entry.Name)
	if err != nil {
		return nil, err
	}
	if err = format.PutSizeInHeader(t.content[:format.HeaderSize], entry.Size); err != nil {
		return nil, err
	}
	t.content = append(t.content, entry.Data...)
	return t, nil
}

func Parse(content []byte) ([]entry.Entry, error) {
	var entries []entry.Entry

	for sepPos := 0; len(content) > 0; content = content[sepPos+1:] {
		if len(content) < format.HeaderSize {
			return nil, errors.New("error parsing tree object: invalid format")
		}

		sepPos = bytes.IndexByte(content[hasher.HashSize+1:], format.TreeEntriesSeparator)
		if sepPos == -1 {
			return nil, errors.New("error parsing tree object: invalid format")
		}
		sepPos += hasher.HashSize + 1

		entries = append(
			entries, entry.Entry{
				Type: content[0],
				Hash: content[1 : hasher.HashSize+1],
				Name: string(content[hasher.HashSize+1 : sepPos]),
			},
		)
	}

	return entries, nil
}

func (t *Tree) RegisterEntry(child entry.Interface) error {
	t.content = append(t.content, child.GetType())
	t.content = append(t.content, child.GetHash()...)
	t.content = append(t.content, child.GetName()...)
	t.content = append(t.content, format.TreeEntriesSeparator)
	return format.PutSizeInHeader(t.content[:format.HeaderSize], len(t.content)-format.HeaderSize)
}

func (t *Tree) Clear(name string) (*Tree, error) {
	newTree, err := Prepare(name)
	if err != nil {
		return nil, err
	}

	entries, err := Parse(t.content)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.Type == format.DocsType || e.Type == format.CodeType {
			if err = newTree.RegisterEntry(e); err != nil {
				return nil, err
			}
		}
	}
	return newTree, nil
}

func (t *Tree) GetType() byte {
	return format.TreeType
}

func (t *Tree) GetName() string {
	return t.name
}

func (t *Tree) SetName(name string) {
	t.name = name
}

func (t *Tree) GetHash() []byte {
	return hasher.Calculate(t.content)
}

func (t *Tree) GetData() []byte {
	return t.content
}

func (t *Tree) GetEntries() ([]entry.Entry, error) {
	return Parse(t.content[format.HeaderSize:])
}
