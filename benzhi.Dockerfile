FROM golang:1.22-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/predictive-maintenance ./cmd/server

FROM debian:bookworm-slim

RUN groupadd --gid 10001 app \
    && useradd --uid 10001 --gid app --create-home --home-dir /app app \
    && mkdir --parents /app/data \
    && chown --recursive app:app /app

WORKDIR /app

COPY --from=builder /out/predictive-maintenance /app/predictive-maintenance

ENV HTTP_ADDR=:8080 \
    DB_PATH=/app/data/predictive-maintenance.db

USER app:app

EXPOSE 8080

ENTRYPOINT ["/app/predictive-maintenance"]
