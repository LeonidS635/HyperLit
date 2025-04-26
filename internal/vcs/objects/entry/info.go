package entry

import (
	"bytes"
	"errors"

	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
)

type Info struct {
	Path string
	Hash []byte
}

func GetInfo(e Interface) Info {
	return Info{
		Path: e.GetPath(),
		Hash: e.GetHash(),
	}
}

func ParseInfo(b []byte) (Info, error) {
	if len(b) < hasher.HashSize {
		return Info{}, errors.New("invalid info")
	}
	return Info{
		Path: string(b[:len(b)-hasher.HashSize]),
		Hash: b[len(b)-hasher.HashSize:],
	}, nil
}

func (i Info) Dump() []byte {
	b := bytes.NewBuffer(nil)
	b.WriteString(i.Path)
	b.Write(i.Hash)
	return b.Bytes()
}
