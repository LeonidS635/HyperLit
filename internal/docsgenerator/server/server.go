package server

import (
	"fmt"
	"net/http"
	"os"
)

func Start(port int, htmlFilePath string) error {
	if _, err := os.Stat(htmlFilePath); err != nil {
		return err
	}

	http.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, htmlFilePath)
		},
	)
	http.HandleFunc("/open-file", openFileHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func openFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "no file specified", http.StatusBadRequest)
		return
	}

	fmt.Println(fileName)
	content, err := os.ReadFile(fileName)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(content)
}
