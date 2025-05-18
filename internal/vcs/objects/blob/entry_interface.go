package blob

import (
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
)

func (b *Blob) GetType() byte {
	return b.type_
}

func (b *Blob) GetName() string {
	return ""
}

func (b *Blob) GetHash() []byte {
	return hasher.Calculate(b.content)
}

func (b *Blob) GetData() []byte {
	return b.content
}

func (b *Blob) GetContent() []byte {
	return b.content[format.HeaderSize:]
}
