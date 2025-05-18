package info

const (
	StatusUnmodified = iota
	StatusProbablyModified
	StatusModified

	StatusDocsOutdated
	StatusCodeOutdated

	StatusCreated
	StatusDeleted
)

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
