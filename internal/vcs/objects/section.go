package objects

import (
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

type Section struct {
	Path string
	Code entry.Entry
	Docs entry.Entry
}
