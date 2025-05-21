package handler

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/go-playground/validator/v10"
    "go.temporal.io/sdk/client"

    "temporal-mcp-gateway/internal/model"
)

type WorkflowStarter interface {
    ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow string, args ...interface{}) (client.WorkflowRun, error)
}

type StartWorkflowRequest struct {
    Namespace  string      `json:"namespace" validate:"required"`
    WorkflowID string      `json:"workflow_id" validate:"required"`
    TaskQueue  string      `json:"task_queue" validate:"required"`
    Input      interface{} `json:"input"`
}

type StartWorkflowResponse struct {
    RunID          string `json:"run_id"`
    AlreadyRunning bool   `json:"already_running"`
}

func StartWorkflow(c WorkflowStarter, v *validator.Validate) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req StartWorkflowRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            model.NewError("bad_request", err.Error()).Write(w, http.StatusBadRequest)
            return
        }
        if err := v.Struct(req); err != nil {
            model.NewError("bad_request", err.Error()).Write(w, http.StatusBadRequest)
            return
        }
        opts := client.StartWorkflowOptions{
            ID:        req.WorkflowID,
            TaskQueue: req.TaskQueue,
        }
        we, err := c.ExecuteWorkflow(r.Context(), opts, "workflow", req.Input)
        if err != nil {
            if client.IsWorkflowExecutionAlreadyStartedError(err) {
                model.NewError("already_started", err.Error()).Write(w, http.StatusConflict)
                return
            }
            model.NewError("temporal_error", err.Error()).Write(w, http.StatusInternalServerError)
            return
        }
        resp := StartWorkflowResponse{RunID: we.GetRunID()}
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(resp)
    }
}
