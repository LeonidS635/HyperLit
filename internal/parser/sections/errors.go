package sections

import (
	"errors"
	"fmt"
)

type ParseErr struct {
	line int
	err  error
}

func (e ParseErr) Error() string {
	return fmt.Sprintf("error parsing sections: line %d: %v", e.line, e.err)
}

var CloseUnopenedDocsErr = errors.New("close unopened docs")
var CloseUnopenedCodeErr = fmt.Errorf("close unopened code")

var DuplicateSectionNameErr = fmt.Errorf("duplicate section name")
var IncorrectSectionNameErr = fmt.Errorf("incorrect section name")
var EmptySectionNameErr = fmt.Errorf("empty section name")
