# Build stage: compile static Go binary
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY mcp/ ./mcp/
COPY mcp/supabase/ ./mcp/supabase/
WORKDIR /app/mcp
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# Final stage: minimal image
FROM scratch
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
