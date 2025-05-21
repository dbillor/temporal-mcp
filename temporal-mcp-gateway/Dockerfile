FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /temporal-mcp-gateway ./cmd/server

FROM gcr.io/distroless/base
WORKDIR /
COPY --from=builder /temporal-mcp-gateway /temporal-mcp-gateway
CMD ["/temporal-mcp-gateway"]
