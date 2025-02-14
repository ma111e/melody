package router

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/ma111e/melody/internal/events"

	"github.com/ma111e/melody/internal/logging"

	"github.com/ma111e/melody/internal/config"
)

// StartHTTP starts the dummy HTTP server
func StartHTTP(quitErrChan chan error) {
	r := http.NewServeMux()
	r.Handle("/",
		headersHandler(
			melodyFs(http.Dir(config.Cfg.ServerHTTPDir), config.Cfg.ServerHTTPMissingResponseStatus),
			config.Cfg.ServerHTTPHeaders))

	logging.Std.Println("Started HTTP server on port", config.Cfg.ServerHTTPPort)
	quitErrChan <- http.ListenAndServe(fmt.Sprintf(":%d", config.Cfg.ServerHTTPPort), r)
}

// StartHTTPS starts the dummy HTTPS server
func StartHTTPS(quitErrChan chan error, eventChan chan events.Event) {
	r := http.NewServeMux()
	r.Handle("/",
		httpsLogger(
			headersHandler(
				melodyFs(http.Dir(config.Cfg.ServerHTTPSDir), config.Cfg.ServerHTTPSMissingResponseStatus),
				config.Cfg.ServerHTTPSHeaders), eventChan))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Cfg.ServerHTTPSPort),
		Handler:      r,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	logging.Std.Println("Started HTTPS server on port", config.Cfg.ServerHTTPSPort)
	quitErrChan <- srv.ListenAndServeTLS(config.Cfg.ServerHTTPSCert, config.Cfg.ServerHTTPSKey)
}
