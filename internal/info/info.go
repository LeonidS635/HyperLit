package info

import (
	"fmt"
	"sync"
	"time"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

type Section struct {
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

const (
	StatusUnmodified = iota
	StatusProbablyModified
	StatusModified

	StatusDocsOutdated
	StatusCodeOutdated

	StatusCreated
	StatusDeleted
)

type SectionStatus struct {
	Path         string
	Trie         *trie.Node[Section]
	FullTrieNode *trie.Node[TrieSection]
}

type TrieSection struct {
	Section *tree.Tree
	Status  int
}

type SectionsStatuses struct {
	mu            sync.Mutex
	pathsByStatus map[int][]SectionStatus
}

func newSectionsStatuses() *SectionsStatuses {
	return &SectionsStatuses{
		pathsByStatus: make(map[int][]SectionStatus),
	}
}

func (s *SectionsStatuses) Check(statusCode int) bool {
	_, ok := s.pathsByStatus[statusCode]
	return ok
}

func (s *SectionsStatuses) Add(statusCode int, status SectionStatus) {
	s.mu.Lock()
	s.pathsByStatus[statusCode] = append(s.pathsByStatus[statusCode], status)
	s.mu.Unlock()
}

func (s *SectionsStatuses) Get(statusCode int) []SectionStatus {
	return s.pathsByStatus[statusCode]
}

func (s *SectionsStatuses) Remove(statusCode int) {
	delete(s.pathsByStatus, statusCode)
}

func (s *SectionsStatuses) Print() {
	for statusCode, statuses := range s.pathsByStatus {
		fmt.Printf("Sections with status %d:\n", statusCode)
		for _, status := range statuses {
			fmt.Printf("\t%s\n", status.Path)
			status.Trie.Print()
			status.FullTrieNode.Print()
		}
	}
}

func compareFileAndSection(fileInfo File, sectionInfo Section) int {
	if fileInfo.MTime.Before(sectionInfo.MTime) {
		return StatusUnmodified
	}
	return StatusProbablyModified
}

func compareTwoSections(newSectionInfo Section, prevSectionInfo Section) int {
	if newSectionInfo.Hash == prevSectionInfo.Hash {
		return StatusUnmodified
	}
	if newSectionInfo.DocsHash != prevSectionInfo.DocsHash && newSectionInfo.CodeHash != prevSectionInfo.CodeHash {
		return StatusModified
	}
	if newSectionInfo.DocsHash == prevSectionInfo.DocsHash && newSectionInfo.CodeHash != prevSectionInfo.CodeHash {
		return StatusDocsOutdated
	}
	if newSectionInfo.CodeHash == prevSectionInfo.CodeHash && newSectionInfo.DocsHash != prevSectionInfo.DocsHash {
		return StatusCodeOutdated
	}
	return StatusUnmodified
}
