package tree

import (
	"sort"

	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/format"
)

type registrations struct {
	isSorted bool

	header   []byte
	data     []byte
	children []entry.Interface
}

func newRegistrations(header []byte) registrations {
	return registrations{
		isSorted: true,
		header:   header,
		data:     make([]byte, len(header)),
	}
}

func (r registrations) append(entry entry.Interface) registrations {
	r.children = append(r.children, entry)
	r.isSorted = false
	return r
}

func (r registrations) setData(data []byte) registrations {
	r.data = make([]byte, len(data))
	copy(r.data, data)
	r.isSorted = true
	return r
}

func (r registrations) getData() ([]byte, error) {
	if r.isSorted {
		return r.data, nil
	}

	sort.Slice(
		r.children, func(i, j int) bool {
			return r.children[i].GetName() < r.children[j].GetName()
		},
	)
	r.isSorted = true

	copy(r.data, r.header)
	for _, child := range r.children {
		r.data = append(r.data, child.GetType())
		r.data = append(r.data, child.GetHash()...)
		r.data = append(r.data, child.GetName()...)
		r.data = append(r.data, format.TreeEntriesSeparator)
	}
	if err := format.PutSizeInHeader(r.data[:format.HeaderSize], len(r.data)-format.HeaderSize); err != nil {
		return nil, err
	}
	return r.data, nil
}
