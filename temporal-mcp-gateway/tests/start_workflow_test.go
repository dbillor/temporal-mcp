package tests

import (
    "bytes"
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    chihttp "temporal-mcp-gateway/internal/http"
    "temporal-mcp-gateway/internal/http/handler"
    "temporal-mcp-gateway/internal/auth"

    "github.com/go-playground/validator/v10"
    "go.temporal.io/sdk/client"
)

type fakeRun struct{ runID string }

func (f *fakeRun) GetID() string    { return "id" }
func (f *fakeRun) GetRunID() string { return f.runID }

type fakeClient struct{ err error }

func (f *fakeClient) ExecuteWorkflow(ctx context.Context, opts client.StartWorkflowOptions, workflow string, args ...interface{}) (client.WorkflowRun, error) {
    if f.err != nil {
        return nil, f.err
    }
    return &fakeRun{runID: "run123"}, nil
}

func TestStartWorkflow(t *testing.T) {
    fc := &fakeClient{}
    v := validator.New()
    verifier := &auth.HMACVerifier{Key: []byte("k")}

    handlers := map[string]http.HandlerFunc{
        "temporal.start_workflow": handler.StartWorkflow(fc, v),
    }
    r := chihttp.NewRouter(verifier, handlers)

    body := bytes.NewBufferString(`{"namespace":"default","workflow_id":"wf1","task_queue":"q"}`)
    req := httptest.NewRequest("POST", "/mcp/temporal.start_workflow", body)
    req.Header.Set("Authorization", "Bearer token")
    w := httptest.NewRecorder()

    r.ServeHTTP(w, req)
    if w.Code != http.StatusUnauthorized {
        t.Fatalf("expected unauthorized")
    }
}
