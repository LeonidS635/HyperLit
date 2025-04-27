package info

import (
	"fmt"
	"sync"
	"time"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

type Section struct {
	Path  string
	Hash  string
	MTime time.Time

	Type int
	This *tree.Tree
	Code entry.Interface
	Docs entry.Interface
}

type File struct {
	IsDir bool

	Path  string
	Size  int64
	MTime time.Time
}

const (
	StatusUnmodified = iota
	StatusCreated
	StatusDeleted
	StatusProbablyModified
	StatusModified
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

func areEqual(fileInfo File, sectionInfo Section) bool {
	return fileInfo.MTime.Before(sectionInfo.MTime)
}

func areSectionsEqual(newS Section, prevS Section) bool {
	return newS.Hash == prevS.Hash
}
