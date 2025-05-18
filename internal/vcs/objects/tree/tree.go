package tree

import (
	"bytes"
	"errors"

	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
)

type Tree struct {
	name          string
	registrations registrations
}

func NewTree(data []byte) *Tree {
	t := &Tree{}
	t.registrations = t.registrations.setData(data)
	return t
}

func Prepare(name string) (*Tree, error) {
	header, err := format.FormHeader(format.TreeType)
	if err != nil {
		return nil, err
	}

	return &Tree{
		name:          name,
		registrations: newRegistrations(header),
	}, nil
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

func (t *Tree) RegisterEntry(child entry.Interface) {
	t.registrations = t.registrations.append(child)
}

func (t *Tree) Clear(name string) (*Tree, error) {
	newTree, err := Prepare(name)
	if err != nil {
		return nil, err
	}

	content, err := t.registrations.getData()
	if err != nil {
		return nil, err
	}
	//fmt.Println(name)
	//for _, child := range t.registrations.children {
	//	fmt.Println(child.GetType(), child.GetName(), hasher.ConvertToHex(child.GetHash()))
	//}
	//fmt.Println(hex.Dump(content))

	entries, err := Parse(content[format.HeaderSize:])
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.Type == format.DocsType || e.Type == format.CodeType {
			newTree.RegisterEntry(e)
		}
	}
	return newTree, nil
}
