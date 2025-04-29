package html

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
	"github.com/LeonidS635/HyperLit/internal/vcs/hasher"
)

const (
	sectionTemplate = "<span class=folder onclick=toggleVisibility('%s')>%s</span>"
	ulOpenTemplate  = "<ul id=%s class='hidden nested'>"
)

func Generate(htmlFilePath string, rootNode *trie.Node[info.TrieSection], rootName string) error {
	htmlFile, err := os.Create(htmlFilePath)
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	writer := bufio.NewWriter(htmlFile)
	defer writer.Flush()

	if _, err = writer.WriteString("<!DOCTYPE html><html lang=ru><head>"); err != nil {
		return err
	}
	if _, err = writer.WriteString(head); err != nil {
		return err
	}
	if _, err = writer.WriteString(css); err != nil {
		return err
	}
	if _, err = writer.WriteString("<body>"); err != nil {
		return err
	}
	if _, err = writer.WriteString("<div class=container><div class=tree><ul><li>"); err != nil {
		return err
	}
	if err = gen(rootNode, rootName, writer); err != nil {
		return err
	}
	if _, err = writer.WriteString("</li></ul></div>" + helloPage + "</div>"); err != nil {
		return err
	}
	if _, err = writer.WriteString(script); err != nil {
		return err
	}
	if _, err = writer.WriteString("</body></html>"); err != nil {
		return err
	}

	return nil
}

func gen(node *trie.Node[info.TrieSection], name string, writer *bufio.Writer) error {
	if node.Data.Status == info.StatusDeleted {
		return nil
	}

	hash := hasher.ConvertToHex(node.Data.Section.GetHash())
	filePath := filepath.Join(hash[:2], hash[2:])

	if _, err := writer.WriteString(fmt.Sprintf(sectionTemplate, filePath, name)); err != nil {
		return err
	}
	if _, err := writer.WriteString(fmt.Sprintf(ulOpenTemplate, filePath)); err != nil {
		return err
	}

	for childName, child := range node.GetAll() {
		if _, err := writer.WriteString("<li>"); err != nil {
			return err
		}
		if err := gen(child, childName, writer); err != nil {
			return err
		}
		if _, err := writer.WriteString("</li>"); err != nil {
			return err
		}
	}

	_, err := writer.WriteString("</ul>")
	return err
}

func FormDocumentation(docs, code []byte) []byte {
	documentation := make([]byte, 0, len(docs)+len(code))

	documentation = append(documentation, docs...)
	documentation = append(documentation, '\n')
	documentation = append(documentation, "<details><summary>Показать код</summary><pre><code>"...)
	documentation = append(documentation, code...)
	documentation = append(documentation, "</code></pre></details>"...)

	return documentation
}
