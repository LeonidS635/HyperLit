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
	path string
}

func NewGenerator(path string) Generator {
	return Generator{path: filepath.Join(path, indexHTMLFile)}
}

func (g Generator) Generate(rootNode *trie.Node[info.Section], rootName string) error {
	return html.Generate(g.path, rootNode, rootName)
}

func (g Generator) StartServer() error {
	return server.Start(serverPort, g.path)
}
