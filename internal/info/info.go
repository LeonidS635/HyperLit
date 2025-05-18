package info

import (
	"time"

	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

type Section struct {
	Status int
	Path   string

	Hash     string
	CodeHash string
	DocsHash string

	MTime time.Time
	This  *tree.Tree
}

type File struct {
	IsDir bool
	Path  string
	Size  int64
	MTime time.Time
}
