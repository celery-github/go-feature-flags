package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/celery-github/go-feature-flags/internal/flags"
)

type Handlers struct {
	svc    *flags.Service
	logger *log.Logger
}

func NewHandlers(svc *flags.Service, logger *log.Logger) *Handlers {
	return &Handlers{svc: svc, logger: logger}
}

func (h *Handlers) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handlers) FlagsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		all, err := h.svc.List()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "list_failed", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"flags": all})
	case http.MethodPost:
		var f flags.Flag
		if err := decodeJSON(r.Body, &f); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		created, err := h.svc.Put(f)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "validation_failed", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)
	default:
		writeErr(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET or POST")
	}
}

func (h *Handlers) FlagsItem(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/flags/")
	if name == "" || strings.Contains(name, "/") {
		writeErr(w, http.StatusBadRequest, "invalid_name", "flag name required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		f, err := h.svc.Get(name)
		if err != nil {
			if errors.Is(err, flags.ErrNotFound) {
				writeErr(w, http.StatusNotFound, "not_found", "flag not found")
				return
			}
			writeErr(w, http.StatusInternalServerError, "get_failed", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, f)
	case http.MethodPatch:
		var patch flags.FlagUpsert
		if err := decodeJSON(r.Body, &patch); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		updated, err := h.svc.Patch(name, patch)
		if err != nil {
			if errors.Is(err, flags.ErrNotFound) {
				writeErr(w, http.StatusNotFound, "not_found", "flag not found")
				return
			}
			writeErr(w, http.StatusBadRequest, "patch_failed", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		if err := h.svc.Delete(name); err != nil {
			if errors.Is(err, flags.ErrNotFound) {
				writeErr(w, http.StatusNotFound, "not_found", "flag not found")
				return
			}
			writeErr(w, http.StatusInternalServerError, "delete_failed", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": name})
	default:
		writeErr(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET, PATCH, or DELETE")
	}
}

func (h *Handlers) Evaluate(w http.ResponseWriter, r *http.Request) {
	// /evaluate/{name}?env=prod&user=celeste
	name := strings.TrimPrefix(r.URL.Path, "/evaluate/")
	if name == "" || strings.Contains(name, "/") {
		writeErr(w, http.StatusBadRequest, "invalid_name", "flag name required")
		return
	}
	env := r.URL.Query().Get("env")
	if env == "" {
		env = "dev"
	}
	userKey := r.URL.Query().Get("user")

	on, f, err := h.svc.Evaluate(name, env, userKey)
	if err != nil {
		if errors.Is(err, flags.ErrNotFound) {
			writeErr(w, http.StatusNotFound, "not_found", "flag not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "evaluate_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"name":        f.Name,
		"env":         env,
		"user":        userKey,
		"enabled":     on,
		"rolloutType": f.Rollout.Type,
		"percentage":  f.Rollout.Percentage,
	})
}

func decodeJSON(body io.ReadCloser, v any) error {
	defer body.Close()
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
