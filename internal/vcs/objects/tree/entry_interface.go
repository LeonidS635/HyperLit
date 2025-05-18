package tree

import (
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
)

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
	content, _ := t.registrations.getData()
	return hasher.Calculate(content)
}

func (t *Tree) GetData() []byte {
	content, _ := t.registrations.getData()
	return content
}

func (t *Tree) GetContent() []byte {
	content, _ := t.registrations.getData()
	return content[format.HeaderSize:]
}
