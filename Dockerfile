# syntax=docker/dockerfile:1.7

# ----- builder -----
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Cache modules first.
COPY go.mod go.sum* ./
RUN go mod download

# Source.
COPY . .

# CGO_ENABLED=0 keeps the binary statically linked for alpine.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/cli    ./cmd/cli

# ----- runner -----
FROM alpine:3.20 AS runner

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder --chown=app:app /out/server /app/server
COPY --from=builder --chown=app:app /out/cli    /app/cli
COPY --from=builder --chown=app:app /src/migrations /app/migrations
COPY --from=builder --chown=app:app /src/config     /app/config

EXPOSE 8080
USER app
ENTRYPOINT ["/app/server"]
