package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

const message = "Hello World! Surya's here!"
const (
	address        = ":8080"
	certificateDir = "./storage/local.cert"
	privateKeyDir  = "./storage/local.pem"
)

type Handlers struct {
	logger  *log.Logger
	db      *sqlx.DB
	context context
}

type context struct {
	req          *http.Request
	res          http.ResponseWriter
	responseTime time.Duration
}

func main() {

	logger := log.New(os.Stdout, "log-", log.LstdFlags|log.Lshortfile)

	db, err := sqlx.Open("postgres", "postgres://postgres@localhost:5432/simple")
	if err != nil {
		logger.Fatalf("Failed to connect database: %s", err)
	}

	handler := newHandlers(logger, db)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Logger(handler.rootHandler))

	server := createNewServer(mux, address)

	if err := server.ListenAndServeTLS(certificateDir, privateKeyDir); err != nil {
		logger.Fatalf("Server failed to start: %s", err)
	}
}

func createNewServer(mux *http.ServeMux, serverAddress string) *http.Server {
	// https://blog.cloudflare.com/exposing-go-on-the-internet/
	tlsConfig := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},

		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	return &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}
}

func newHandlers(logger *log.Logger, db *sqlx.DB) *Handlers {
	return &Handlers{
		logger: logger,
		db:     db,
	}
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("Response time: %s\n", time.Now().Sub(startTime))

		next(w, r)
	}
}

func (h *Handlers) rootHandler(w http.ResponseWriter, r *http.Request) {
	h.context.req = r
	h.context.res = w

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (h *Handlers) getHandler(w http.ResponseWriter, r *http.Request) {
	h.context.req = r
	h.context.res = w

	h.db.ExecContext(r.Context)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
