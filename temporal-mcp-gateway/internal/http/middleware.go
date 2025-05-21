package http

import (
    "context"
    "log"
    "net/http"
    "strings"

    "temporal-mcp-gateway/internal/auth"
    "temporal-mcp-gateway/internal/model"
)

// JSONError wraps handlers and writes errors in MCP format.
func JSONError(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rw := &responseWriter{ResponseWriter: w}
        next.ServeHTTP(rw, r)
        if rw.err != nil {
            model.NewError("internal_error", rw.err.Error()).Write(w, rw.status)
        }
    })
}

type responseWriter struct {
    http.ResponseWriter
    status int
    err    error
}

func (w *responseWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(b []byte) (int, error) {
    n, err := w.ResponseWriter.Write(b)
    w.err = err
    return n, err
}

// Logger prints a basic request log line.
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// Auth verifies JWT tokens and scope for the given tool.
func Auth(verifier auth.Verifier, tool string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authz := r.Header.Get("Authorization")
            token := strings.TrimPrefix(authz, "Bearer ")
            if token == "" {
                model.NewError("unauthorized", "missing token").Write(w, http.StatusUnauthorized)
                return
            }
            claims, extra, err := verifier.Verify(token)
            if err != nil {
                model.NewError("unauthorized", err.Error()).Write(w, http.StatusUnauthorized)
                return
            }
            if scopes, ok := extra["scp"].([]interface{}); ok {
                allowed := false
                for _, s := range scopes {
                    if str, ok := s.(string); ok && str == tool {
                        allowed = true
                        break
                    }
                }
                if !allowed {
                    model.NewError("forbidden", "scope not allowed").Write(w, http.StatusForbidden)
                    return
                }
            }
            ctx := context.WithValue(r.Context(), "sub", claims.Subject)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
