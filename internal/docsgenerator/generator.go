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

	getDataByHash func(hash string) ([]byte, error)
}

func NewGenerator(rootPath string, getDataByHash func(hash string) ([]byte, error)) Generator {
	return Generator{
		htmlFilepath:  filepath.Join(rootPath, indexHTMLFile),
		getDataByHash: getDataByHash,
	}
}

func (g Generator) Generate(rootNode *trie.Node[info.Section], rootName string) error {
	return html.Generate(g.htmlFilepath, rootNode, rootName)
}

func (g Generator) StartServer() error {
	return server.Start(serverPort, g.htmlFilepath, g.getDataByHash)
}
