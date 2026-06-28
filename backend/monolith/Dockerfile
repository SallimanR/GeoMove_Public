FROM golang:1.26-alpine AS base

WORKDIR /app

ENV GOCACHE=/root/.cache/go-build
ENV	GOMODCACHE=/go/pkg/mod
ENV GOBIN=/go/bin

COPY go.mod go.sum ./

RUN --mount=type=cache,target=${GOCACHE} \
	--mount=type=cache,target=${GOMODCACHE} \
	go mod download

FROM base AS dev

RUN --mount=type=cache,target=${GOCACHE} \
	--mount=type=cache,target=${GOMODCACHE} \
	go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]


FROM golang:1.26-alpine AS build

ENV GOCACHE=/root/.cache/go-build \
	GOMODCACHE=/go/pkg/mod \
	GOBIN=/go/bin

COPY . .

RUN --mount=type=cache,target=${GOCACHE} \
	--mount=type=cache,target=${GOMODCACHE} \
	go build -o ./main ./cmd/main.go

CMD ["./main"]
