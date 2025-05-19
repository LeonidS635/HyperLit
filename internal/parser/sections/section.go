package sections

import (
	"bufio"
	"time"

	"github.com/LeonidS635/HyperLit/internal/parser/comments"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
)

var (
	docsStartSeq = []byte("@@docs")
	docsEndSeq   = []byte("@@/docs")
	codeEndSeq   = []byte("@@/code")
)

type Parser struct {
	commentsAnalyzer *comments.Analyzer
	fileScanner      *bufio.Scanner
	fileModTime      time.Time
	blobsSavingCh    chan<- entry.Interface
	lineNumber       int
}

func NewParser(
	filePath string, fileModTime time.Time, fileScanner *bufio.Scanner, blobsSavingCh chan<- entry.Interface,
) (*Parser, error) {
	commentsAnalyzer, err := comments.NewAnalyzer(filePath)
	if err != nil {
		return nil, err
	}

	return &Parser{
		commentsAnalyzer: commentsAnalyzer,
		fileScanner:      fileScanner,
		fileModTime:      fileModTime,
		blobsSavingCh:    blobsSavingCh,
		lineNumber:       0,
	}, nil
}
