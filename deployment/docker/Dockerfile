FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=${APP_VERSION:-1}" \
    -o /app/main \
    ./cmd/app/main.go

FROM gcr.io/distroless/static-debian12

ENV APP_NAME="TEST_APP" \
    HTTP_PORT="8000" \
    LOG_LEVEL="info" \
    PG_POOL_MAX="10" \
    PG_POOL_IDLE="1"

EXPOSE 8000
USER nonroot:nonroot
COPY --from=builder /app/main /
ENTRYPOINT ["/main"]
