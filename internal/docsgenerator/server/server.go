package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func Start(ctx context.Context, port int, htmlFilePath string, getDataByHash func(hash string) ([]byte, error)) error {
	if _, err := os.Stat(htmlFilePath); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, htmlFilePath)
		},
	)
	mux.HandleFunc(
		"/gen", func(w http.ResponseWriter, r *http.Request) {
			openFileHandler(w, r, getDataByHash)
		},
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
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

	documentation, err := json.Marshal(
		struct {
			Docs string `json:"docs"`
			Code string `json:"code"`
		}{
			Docs: string(docs),
			Code: string(code),
		},
	)
	if err != nil {
		http.Error(w, "error marshalling documentation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(documentation)
}
