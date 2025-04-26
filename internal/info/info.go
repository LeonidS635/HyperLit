package info

import "time"

type Section struct {
	MTime time.Time
}

type File struct {
	IsDir bool

	MTime time.Time
}

const (
	StatusUnmodified = iota
	StatusCreated
	StatusDeleted
	StatusRenamed
	StatusModified
)

type FileStatus struct {
	Path   string
	Status int
}

func areEqual(fileInfo File, sectionInfo Section) bool {
	return true
}
