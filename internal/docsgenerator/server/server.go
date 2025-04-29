package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LeonidS635/HyperLit/internal/docsgenerator/html"
)

func Start(port int, htmlFilePath string, parseFileFn func(hash string) ([]byte, []byte, error)) error {
	if _, err := os.Stat(htmlFilePath); err != nil {
		return err
	}

	http.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, htmlFilePath)
		},
	)
	http.HandleFunc(
		"/open-file", func(w http.ResponseWriter, r *http.Request) {
			openFileHandler(w, r, parseFileFn)
		},
	)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func openFileHandler(w http.ResponseWriter, r *http.Request, parseFileFn func(hash string) ([]byte, []byte, error)) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "no file specified", http.StatusBadRequest)
		return
	}

	docs, code, err := parseFileFn(fileName)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(html.FormDocumentation(docs, code))
}
