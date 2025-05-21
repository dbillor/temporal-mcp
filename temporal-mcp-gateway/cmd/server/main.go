package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    chihttp "temporal-mcp-gateway/internal/http"
    "temporal-mcp-gateway/internal/http/handler"
    "temporal-mcp-gateway/internal/auth"
    "temporal-mcp-gateway/internal/temporal"

    "github.com/go-playground/validator/v10"
)

func main() {
    addr := os.Getenv("TEMPORAL_ADDRESS")
    if addr == "" {
        addr = "localhost:7233"
    }

    key := os.Getenv("JWT_SIGNING_KEY")
    verifier := &auth.HMACVerifier{Key: []byte(key)}

    tc, err := temporal.NewClient(addr)
    if err != nil {
        log.Fatal(err)
    }
    defer temporal.CloseSafe(tc)

    v := validator.New()

    handlers := map[string]http.HandlerFunc{
        "temporal.start_workflow": handler.StartWorkflow(tc, v),
    }

    r := chihttp.NewRouter(verifier, handlers)
    srv := &http.Server{Addr: ":8080", Handler: r}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _ = srv.Shutdown(ctx)
}
