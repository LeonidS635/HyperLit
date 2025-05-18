package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LeonidS635/HyperLit/internal/docsgenerator/html"
)

func Start(port int, htmlFilePath string, getDataByHash func(hash string) ([]byte, error)) error {
	if _, err := os.Stat(htmlFilePath); err != nil {
		return err
	}

	http.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, htmlFilePath)
		},
	)
	http.HandleFunc(
		"/gen", func(w http.ResponseWriter, r *http.Request) {
			openFileHandler(w, r, getDataByHash)
		},
	)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func openFileHandler(w http.ResponseWriter, r *http.Request, getDataByHashFn func(hash string) ([]byte, error)) {
	codeHash := r.URL.Query().Get("code")
	docsHash := r.URL.Query().Get("docs")

	var code, docs []byte
	var err error

	if codeHash != "" {
		code, err = getDataByHashFn(codeHash)
		if err != nil {
			http.Error(w, fmt.Sprintf("code file %s not found", codeHash), http.StatusNotFound)
			return
		}
	}
	if docsHash != "" {
		docs, err = getDataByHashFn(docsHash)
		if err != nil {
			http.Error(w, fmt.Sprintf("docs file %s not found", docsHash), http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(html.FormDocumentation(docs, code))
}
