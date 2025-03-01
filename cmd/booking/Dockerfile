FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/service ./cmd/booking

# --- Final Image ---
FROM gcr.io/distroless/static

WORKDIR /go/bin
COPY ./cmd/booking/config/ ./config/
COPY --from=builder /bin/service .

ENTRYPOINT ["./service"]