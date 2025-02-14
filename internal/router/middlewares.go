package router

import (
	"net/http"
	"path/filepath"

	"github.com/ma111e/melody/internal/config"

	"github.com/ma111e/melody/internal/logging"

	"github.com/ma111e/melody/internal/events"
)

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func headersHandler(h http.Handler, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for header, val := range headers {
			w.Header().Set(header, val)
		}
		h.ServeHTTP(w, r) // pass request
	})
}

func httpsLogger(h http.Handler, eventChan chan events.Event) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := config.Cfg.DiscardProto4[config.HTTPSKind]; ok {
			h.ServeHTTP(w, r) // pass request
			return
		} else if _, ok := config.Cfg.DiscardProto6[config.HTTPSKind]; ok {
			h.ServeHTTP(w, r) // pass request
			return
		}

		ev, err := events.NewHTTPEventFromRequest(r)
		if err != nil {
			logging.Errors.Println(err)
			return
		}
		eventChan <- ev

		h.ServeHTTP(w, r) // pass request
	})
}
