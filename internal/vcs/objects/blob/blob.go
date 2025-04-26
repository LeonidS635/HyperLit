package blob

import (
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
)

type Blob struct {
	type_   byte
	content []byte
}

func prepare(type_ byte) (*Blob, error) {
	header, err := format.FormHeader(type_)
	if err != nil {
		return nil, err
	}

	return &Blob{
		type_:   type_,
		content: header,
	}, nil
}

func PrepareCode() (*Blob, error) {
	return prepare(format.CodeType)
}

func PrepareDocs() (*Blob, error) {
	return prepare(format.DocsType)
}
