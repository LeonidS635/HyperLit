package html

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/LeonidS635/HyperLit/internal/docsgenerator/html/static"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
)

const (
	spanTemplate = "<span class=\"folder\" data-code=\"%s\" data-docs=\"%s\">%s</span>"
	ulTemplate   = "<ul class=\"hidden nested\">"
)

func Generate(htmlFilePath string, rootNode *trie.Node[info.Section], rootName string) error {
	htmlFile, err := os.Create(htmlFilePath)
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	writer := bufio.NewWriter(htmlFile)
	defer writer.Flush()

	if _, err = writer.WriteString("<!DOCTYPE html><html lang=ru>"); err != nil {
		return err
	}
	if _, err = writer.WriteString(static.Head); err != nil {
		return err
	}
	if _, err = writer.WriteString(static.Style); err != nil {
		return err
	}
	if _, err = writer.WriteString("<body><div class=\"container\"><div class=\"tree\"><div class=\"tree-header\">HyperLit</div><ul>"); err != nil {
		return err
	}
	if err = gen(rootNode, rootName, writer); err != nil {
		return err
	}
	if _, err = writer.WriteString("</ul></div>" + static.EmptyState + "</div>"); err != nil {
		return err
	}
	if _, err = writer.WriteString(static.Script); err != nil {
		return err
	}
	if _, err = writer.WriteString("</body></html>"); err != nil {
		return err
	}

	return nil
}

func gen(node *trie.Node[info.Section], name string, writer *bufio.Writer) error {
	if node.Data.Status == info.StatusDeleted {
		return nil
	}

	if _, err := writer.WriteString("<li>"); err != nil {
		return err
	}
	if _, err := writer.WriteString(
		fmt.Sprintf(
			spanTemplate, node.Data.CodeHash, node.Data.DocsHash, name,
		),
	); err != nil {
		return err
	}
	if _, err := writer.WriteString(ulTemplate); err != nil {
		return err
	}

	children := node.GetAll()
	childrenNames := make([]string, 0, len(children))
	for childName := range children {
		childrenNames = append(childrenNames, childName)
	}
	sort.Strings(childrenNames)

	for _, childName := range childrenNames {
		if err := gen(children[childName], childName, writer); err != nil {
			return err
		}
	}

	_, err := writer.WriteString("</ul></li>")
	return err
}
