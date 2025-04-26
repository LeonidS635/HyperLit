package blob

import "github.com/LeonidS635/HyperLit/internal/vcs/objects/format"

func (b *Blob) Write(data []byte) error {
	b.content = append(b.content, data...)
	return format.PutSizeInHeader(b.content[:format.HeaderSize], len(b.content)-format.HeaderSize)
}
