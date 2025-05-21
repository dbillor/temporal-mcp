# temporal-mcp

This repository contains an MCP (Model Completion Protocol) specification
for several Temporal workflow management tools. The API definitions are
in `openapi.yaml` and provide the following operations:

- `temporal.list_workflows`
- `temporal.describe_workflow`
- `temporal.start_workflow`
- `temporal.get_workflow_result`
- `temporal.signal_workflow`
- `temporal.query_workflow`
- `temporal.cancel_workflow`
- `temporal.terminate_workflow`

These endpoints can be used by compatible agents to interact with a
Temporal server.
An accompanying `ai-plugin.json` manifest references this OpenAPI file.
