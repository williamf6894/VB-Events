FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/vb-events ./cmd && \
    CGO_ENABLED=0 GOOS=linux go build -o /out/vb-events-seed ./cmd/seed

FROM alpine:3.22

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 vbevents

USER vbevents

COPY --from=builder /out/vb-events /usr/local/bin/vb-events
COPY --from=builder /out/vb-events-seed /usr/local/bin/vb-events-seed

EXPOSE 8200

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8200/healthz || exit 1

CMD ["vb-events"]
