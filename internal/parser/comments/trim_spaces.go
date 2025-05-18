package comments

import "bytes"

const TabWidth = 4

func TrimAndCountLeadingSpaces(line []byte) ([]byte, int) {
	spacesCount := 0
	for _, ch := range line {
		if ch == ' ' {
			spacesCount++
		} else if ch == '\t' {
			spacesCount += TabWidth
		} else {
			break
		}
	}

	return bytes.TrimSpace(line), spacesCount
}
