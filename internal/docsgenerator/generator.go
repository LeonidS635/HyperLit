package docsgenerator

import (
	"path/filepath"

	"github.com/LeonidS635/HyperLit/internal/docsgenerator/html"
	"github.com/LeonidS635/HyperLit/internal/docsgenerator/server"
	"github.com/LeonidS635/HyperLit/internal/helpers/trie"
	"github.com/LeonidS635/HyperLit/internal/info"
)

const indexHTMLFile = "index.html"
const serverPort = 8123

type Generator struct {
	htmlFilepath string

	parseFileFn func(hash string) ([]byte, []byte, error)
}

func NewGenerator(rootPath string, parseFileFn func(hash string) ([]byte, []byte, error)) Generator {
	return Generator{
		htmlFilepath: filepath.Join(rootPath, indexHTMLFile),
		parseFileFn:  parseFileFn,
	}
}

func (g Generator) Generate(rootNode *trie.Node[info.TrieSection], rootName string) error {
	return html.Generate(g.htmlFilepath, rootNode, rootName)
}

func (g Generator) StartServer() error {
	return server.Start(serverPort, g.htmlFilepath, g.parseFileFn)
}
