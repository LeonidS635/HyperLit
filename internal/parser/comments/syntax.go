package comments

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed config/comments.json
var commentsSyntaxJSON []byte

type syntaxJSON struct {
	SingleLine     string `json:"single_line"`
	MultiLineStart string `json:"multi_line_start"`
	MultiLineEnd   string `json:"multi_line_end"`
}

// Syntax struct

type Syntax struct {
	SingleLine     []byte
	MultiLineStart []byte
	MultiLineEnd   []byte
}

var commentsSyntax = make(map[string]Syntax)

// Unmarshall JSON
func init() {
	var syntaxMap map[string]syntaxJSON
	if err := json.Unmarshal(commentsSyntaxJSON, &syntaxMap); err != nil {
		fmt.Println("error parsing comments syntax json:", err)
		os.Exit(1)
	}

	for k, v := range syntaxMap {
		commentsSyntax[k] = Syntax{
			SingleLine:     []byte(v.SingleLine),
			MultiLineStart: []byte(v.MultiLineStart),
			MultiLineEnd:   []byte(v.MultiLineEnd),
		}
	}
}
