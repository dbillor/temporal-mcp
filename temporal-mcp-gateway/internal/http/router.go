package http

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "temporal-mcp-gateway/internal/auth"
)

func NewRouter(verifier auth.Verifier, handlers map[string]http.HandlerFunc) http.Handler {
    r := chi.NewRouter()
    r.Use(Logger)
    r.Use(JSONError)

    for name, h := range handlers {
        r.Method("POST", "/mcp/"+name, Auth(verifier, name)(h))
    }
    return r
}
