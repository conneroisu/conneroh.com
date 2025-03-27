
ARG GO_VERSION=1.24.1
FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .


FROM debian:bookworm

COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
