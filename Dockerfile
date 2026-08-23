FROM golang:1.22-bookworm AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/predictive-maintenance ./cmd/server
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/predictive-maintenance /app/predictive-maintenance
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/predictive-maintenance"]
