package sections

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/parser/comments"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/blob"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/entry"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (p *Parser) Parse(
	ctx context.Context, path string, section *tree.Tree, sectionsTrieNode *trie.Node[info.Section],
) error {
	_, err := p.parse(ctx, path, section, sectionsTrieNode, -1)
	return err
}

func (p *Parser) parse(
	ctx context.Context, path string, section *tree.Tree, sectionsTrieNode *trie.Node[info.Section], sectionOffset int,
) (bool, error) {
	if helpers.IsCtxCancelled(ctx) {
		return false, nil
	}

	sectionsNamesSet := make(map[string]struct{})

	codeSection, err := blob.PrepareCode()
	if err != nil {
		return false, err
	}
	docsSection, err := blob.PrepareDocs()
	if err != nil {
		return false, err
	}

	needToScan := true    // flag if the next line needs to be scanned [false if the recursive call returns due to section offset decreasing, true otherwise]
	isDocsSection := true // flag if we are parsing docs section [true] or code section [false]
	for {
		if needToScan {
			if !p.fileScanner.Scan() {
				break
			}
			p.lineNumber++
		} else {
			needToScan = true
		}

		line := p.fileScanner.Bytes()

		line, lineOffset := comments.TrimAndCountLeadingSpaces(line)
		if len(line) == 0 { // Empty line
			lineOffset = sectionOffset
		}

		// Stop parsing the current section if the offset decreases
		if lineOffset < sectionOffset {
			needToScan = false
			break
		}

		// Check whether the line is a comment line
		lineWithoutCommentSyntax, isComment := p.commentsAnalyzer.IsComment(line)
		isDocsSection = isDocsSection && isComment

		if isComment {
			if docsStartOffset := bytes.Index(lineWithoutCommentSyntax, docsStartSeq); docsStartOffset != -1 {
				//if lineWithoutSpaces := bytes.TrimSpace(line[:docsStartOffset]); len(lineWithoutSpaces) > 0 {
				//	return nil, fmt.Errorf("line %d: %w", *lineNumber, ErrNewLineCmd)
				//}

				// If a new section starts at the same offset, it is a child of the parent section, not ours
				if lineOffset == sectionOffset {
					needToScan = false
					break
				}

				// Check correctness of section name
				name := string(bytes.TrimSpace(lineWithoutCommentSyntax[docsStartOffset+len(docsStartSeq)+1:]))
				if len(name) == 0 { // Section name must not contain only spaces or be empty
					return false, ParseErr{line: p.lineNumber, err: EmptySectionNameErr}
				}
				if _, ok := sectionsNamesSet[name]; ok { // Subsection name must be unique within a single parent section
					return false, ParseErr{line: p.lineNumber, err: DuplicateSectionNameErr}
				}
				sectionsNamesSet[name] = struct{}{}

				// Maintain recursive structure
				subSection, err := tree.Prepare(name)
				if err != nil {
					return false, ParseErr{line: p.lineNumber, err: err}
				}
				nextTrieNode := sectionsTrieNode.Insert(name)

				// Parse child section
				if needToScan, err = p.parse(
					ctx, filepath.Join(path, name), subSection, nextTrieNode, lineOffset,
				); err != nil {
					return false, err
				}
				section.RegisterEntry(subSection)

				// Write section start comment if possible
				var sectionLine []byte
				if singleLineComment := p.commentsAnalyzer.GetSyntax().SingleLine; singleLineComment != nil {
					sectionLine = make([]byte, 0, lineOffset+len(singleLineComment)+len(" Section: ")+len(name))
					for i := 0; i < lineOffset-max(sectionOffset, 0); i++ {
						sectionLine = append(sectionLine, ' ')
					}
					sectionLine = append(sectionLine, singleLineComment...)
					sectionLine = append(sectionLine, ' ')
					sectionLine = append(sectionLine, "Section: "...)
					sectionLine = append(sectionLine, name...)
				} else if multiLineCommentStart, multiLineCommentEnd := p.commentsAnalyzer.GetSyntax().MultiLineStart, p.commentsAnalyzer.GetSyntax().MultiLineEnd; multiLineCommentStart != nil && multiLineCommentEnd != nil {
					sectionLine = make(
						[]byte, 0,
						lineOffset+len(multiLineCommentStart)+len(" Section: ")+len(name)+1+len(multiLineCommentEnd),
					)
					for i := 0; i < lineOffset-max(sectionOffset, 0); i++ {
						sectionLine = append(sectionLine, ' ')
					}
					sectionLine = append(sectionLine, multiLineCommentStart...)
					sectionLine = append(sectionLine, ' ')
					sectionLine = append(sectionLine, "Section: "...)
					sectionLine = append(sectionLine, name...)
					sectionLine = append(sectionLine, ' ')
					sectionLine = append(sectionLine, multiLineCommentEnd...)
				}

				if sectionLine != nil {
					if err = codeSection.WriteLine(sectionLine); err != nil {
						return false, ParseErr{
							line: p.lineNumber,
							err:  fmt.Errorf("error saving secoins line: %w", err),
						}
					}
				}

				continue // Continue parsing from the next line

			} else if docsEndOffset := bytes.Index(lineWithoutCommentSyntax, docsEndSeq); docsEndOffset != -1 {
				//if lineWithoutSpaces := bytes.TrimSpace(line[:docsEndOffset]); len(lineWithoutSpaces) > 0 {
				//	return nil, fmt.Errorf("line %d: %w", *lineNumber, ErrNewLineCmd)
				//}

				// Check that docs section is opened
				if !isDocsSection {
					return false, ParseErr{line: p.lineNumber, err: CloseUnopenedDocsErr}
				}
				isDocsSection = false
				continue // Continue parsing from the next line

			} else if codeEndOffset := bytes.Index(lineWithoutCommentSyntax, codeEndSeq); codeEndOffset != -1 {
				//if lineWithoutSpaces := bytes.TrimSpace(line[:docsEndOffset]); len(lineWithoutSpaces) > 0 {
				//	return nil, fmt.Errorf("line %d: %w", *lineNumber, ErrNewLineCmd)
				//}

				// Check that code section is opened
				if isDocsSection {
					return false, ParseErr{line: p.lineNumber, err: CloseUnopenedCodeErr}
				}
				break // Stop filling current section
			}
		}

		// Write content
		if isDocsSection {
			lineToWrite := make([]byte, 0, lineOffset-sectionOffset+len(lineWithoutCommentSyntax))
			for i := 0; i < lineOffset-max(sectionOffset, 0); i++ {
				lineToWrite = append(lineToWrite, ' ')
			}
			lineToWrite = append(lineToWrite, lineWithoutCommentSyntax...)

			if err = docsSection.WriteLine(lineToWrite); err != nil {
				return false, ParseErr{line: p.lineNumber, err: err}
			}
		} else {
			lineToWrite := make([]byte, 0, lineOffset-sectionOffset+len(line))
			for i := 0; i < lineOffset-max(sectionOffset, 0); i++ {
				lineToWrite = append(lineToWrite, ' ')
			}
			lineToWrite = append(lineToWrite, line...)

			if err = codeSection.WriteLine(lineToWrite); err != nil {
				return false, ParseErr{line: p.lineNumber, err: err}
			}
		}
	}
	if err = p.fileScanner.Err(); err != nil {
		return false, ParseErr{line: p.lineNumber, err: fmt.Errorf("error scanning file: %w", err)}
	}

	// Register docs and code in current section
	section.RegisterEntry(codeSection)
	section.RegisterEntry(docsSection)

	// Fill data in global sections trie
	sectionsTrieNode.Data = info.Section{
		Path: path,

		Hash:     hasher.ConvertToHex(section.GetHash()),
		CodeHash: hasher.ConvertToHex(codeSection.GetHash()),
		DocsHash: hasher.ConvertToHex(docsSection.GetHash()),

		MTime: p.fileModTime,
		This:  section,
	}

	// Send docs and code to save
	helpers.SendCtx[entry.Interface](ctx, p.blobsSavingCh, codeSection)
	helpers.SendCtx[entry.Interface](ctx, p.blobsSavingCh, docsSection)

	return needToScan, nil
}
