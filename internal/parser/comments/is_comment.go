package comments

import "bytes"

// IsComment checks whether a line is a comment line
// (starts with single-line comment characters or is inside a multi-line comment).
//
// If it is a comment line, IsComment returns the line without comment characters and true.
// Otherwise, it returns the unchanged line and false.
func (a Analyzer) IsComment(line []byte) ([]byte, bool) {
	// Multiline case
	if a.isInMultiLineSection {
		// Check if multiline comments section ends
		if bytes.HasPrefix(line, a.syntax.MultiLineEnd) {
			line = bytes.TrimPrefix(line, a.syntax.MultiLineEnd)
			a.isInMultiLineSection = false
		}
		return line, true
	}

	// Single line comment
	if bytes.HasPrefix(line, a.syntax.SingleLine) {
		return bytes.TrimPrefix(line, a.syntax.SingleLine), true
	}

	// Multiline comment
	if bytes.HasPrefix(line, a.syntax.MultiLineStart) {
		a.isInMultiLineSection = true
		return bytes.TrimPrefix(line, a.syntax.MultiLineStart), true
	}

	return line, false
}
