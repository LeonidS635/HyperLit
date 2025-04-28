package parser

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/LeonidS635/HyperLit/internal/helpers"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/blob"
	"github.com/LeonidS635/HyperLit/internal/vcs/objects/tree"
)

func (p *Parser) parseSection(
	ctx context.Context, fileScanner *bufio.Scanner, lineNumber *int, section *tree.Tree,
	curNode *trie.Node[info.Section],
) error {
	if helpers.IsCtxCancelled(ctx) {
		return nil
	}

	codeSection, err := blob.PrepareCode()
	if err != nil {
		return err
	}
	docsSection, err := blob.PrepareDocs()
	if err != nil {
		return err
	}

	isDocsSection := true // flag if we are parsing docs section [true] or code section [false]
	for fileScanner.Scan() {
		line := fileScanner.Bytes()
		*lineNumber++

		if docsStartOffset := bytes.Index(line, p.syntax.DocsStartSeq); docsStartOffset != -1 {
			//if lineWithoutSpaces := bytes.TrimSpace(line[:docsStartOffset]); len(lineWithoutSpaces) > 0 {
			//	return nil, fmt.Errorf("line %d: %w", *lineNumber, ErrNewLineCmd)
			//}

			name := string(bytes.TrimSpace(line[docsStartOffset+len(p.syntax.DocsStartSeq)+1:]))
			subSection, err := tree.Prepare(name)
			if err != nil {
				return err
			}
			nextNode := curNode.Insert(name)

			if err = p.parseSection(ctx, fileScanner, lineNumber, subSection, nextNode); err != nil {
				return err
			}
			if err = section.RegisterEntry(subSection); err != nil {
				return err
			}

		} else if docsEndOffset := bytes.Index(line, p.syntax.DocsEndSeq); docsEndOffset != -1 {
			//if lineWithoutSpaces := bytes.TrimSpace(line[:docsEndOffset]); len(lineWithoutSpaces) > 0 {
			//	return nil, fmt.Errorf("line %d: %w", *lineNumber, ErrNewLineCmd)
			//}

			if !isDocsSection {
				return fmt.Errorf("line %d: %s", *lineNumber, "ErrUnopenedSection")
			}
			isDocsSection = false

		} else if codeEndOffset := bytes.Index(line, p.syntax.CodeEndSeq); codeEndOffset != -1 {
			//if lineWithoutSpaces := bytes.TrimSpace(line[:docsEndOffset]); len(lineWithoutSpaces) > 0 {
			//	return nil, fmt.Errorf("line %d: %w", *lineNumber, ErrNewLineCmd)
			//}

			if isDocsSection {
				return fmt.Errorf("line %d: %s", *lineNumber, "ErrUnopenedSection")
			}
			break

		} else {
			if section == nil {
				return fmt.Errorf("line %d: %s", *lineNumber, "ErrUnopenedSection")
			}

			if isDocsSection {
				docsSection.Write(line)
				docsSection.Write([]byte("\n"))
			} else {
				codeSection.Write(line)
				codeSection.Write([]byte("\n"))
			}
		}
	}
	if err := fileScanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	section.RegisterEntry(codeSection)
	section.RegisterEntry(docsSection)

	curNode.Data = info.Section{
		Hash:     hasher.ConvertToHex(section.GetHash()),
		CodeHash: hasher.ConvertToHex(codeSection.GetHash()),
		DocsHash: hasher.ConvertToHex(docsSection.GetHash()),

		MTime: time.Now(),
		This:  section,
	}

	//docs := docsSection.GetData()
	//code := codeSection.GetData()
	//fmt.Printf(
	//	"Section %q (hash=%s):\n%q\n%q\n%q\n------------------------------------------------------------------\n",
	//	section.GetName(), hasher.ConvertToHex(section.GetHash()),
	//	string(section.GetData()), string(docs), string(code),
	//)
	select {
	case <-ctx.Done():
		return nil
	case p.sectionsCh <- codeSection:
		// ok
	}
	select {
	case <-ctx.Done():
		return nil
	case p.sectionsCh <- docsSection:
		// ok
	}
	select {
	case <-ctx.Done():
		return nil
	case p.sectionsCh <- section:
		// ok
	}

	return nil
}
