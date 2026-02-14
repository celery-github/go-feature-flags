package api

import (
	"log"
	"net/http"

	"github.com/celery-github/go-feature-flags/internal/flags"
)

func NewRouter(svc *flags.Service, logger *log.Logger) http.Handler {
	h := NewHandlers(svc, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.Healthz)

	// Flags CRUD
	mux.Handle("/flags", chain(
		http.HandlerFunc(h.FlagsCollection),
		withJSON,
		withRequestID,
		withLogging(logger),
	))

	mux.Handle("/flags/", chain(
		http.HandlerFunc(h.FlagsItem),
		withJSON,
		withRequestID,
		withLogging(logger),
	))

	// Evaluate
	mux.Handle("/evaluate/", chain(
		http.HandlerFunc(h.Evaluate),
		withJSON,
		withRequestID,
		withLogging(logger),
	))

	return mux
}
