FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o osrs-mcp ./cmd/osrs-mcp

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/osrs-mcp /osrs-mcp
EXPOSE 8080
ENTRYPOINT ["/osrs-mcp", "--addr", ":8080"]
