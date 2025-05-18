package info

import (
	"fmt"
	"sync"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
)

type SectionState struct {
	Path            string
	ProjectTrieNode *trie.Node[Section]
	CurTrie         *trie.Node[Section]
}

type SectionsStates struct {
	mu             sync.Mutex
	statesByStatus map[int][]SectionState
}

func NewSectionsStates() *SectionsStates {
	return &SectionsStates{
		statesByStatus: make(map[int][]SectionState),
	}
}

func (s *SectionsStates) Check(statusCode int) bool {
	_, ok := s.statesByStatus[statusCode]
	return ok
}

func (s *SectionsStates) Add(statusCode int, path string, node, curTrie *trie.Node[Section]) {
	s.mu.Lock()
	s.statesByStatus[statusCode] = append(
		s.statesByStatus[statusCode], SectionState{Path: path, ProjectTrieNode: node, CurTrie: curTrie},
	)
	s.mu.Unlock()
}

func (s *SectionsStates) Get(statusCode int) []SectionState {
	return s.statesByStatus[statusCode]
}

func (s *SectionsStates) Remove(statusCode int) {
	delete(s.statesByStatus, statusCode)
}

func (s *SectionsStates) Print() {
	for statusCode, statuses := range s.statesByStatus {
		fmt.Printf("Sections with status %d:\n", statusCode)
		for _, status := range statuses {
			fmt.Printf("\t%v\n", status)
			status.ProjectTrieNode.Print()
			fmt.Println("children", status.ProjectTrieNode.GetAll())
		}
	}
}
