# Temporal MCP Gateway

This service exposes a set of Temporal operations via the Model Context Protocol (MCP).

## Quick Start

```bash
# start Temporal and the gateway
$ docker-compose up --build
```

The gateway listens on `:8080` by default.

## Environment Variables
- `TEMPORAL_ADDRESS` – Temporal frontend address (default `localhost:7233`)
- `JWT_SIGNING_KEY` – HMAC key used for JWT verification

## Example

```bash
# list namespaces
auth="Authorization: Bearer <token>"
curl -H "$auth" -X POST http://localhost:8080/mcp/temporal.list_namespaces -d '{}'
```


