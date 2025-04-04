ARG         GOLANG_TAG=latest

FROM        golang:${GOLANG_TAG} AS builder

WORKDIR     /api

COPY        go.mod go.sum ./

RUN         go mod download

COPY        . .

RUN         CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o build/api cmd/api/main.go

FROM        scratch

COPY        --from=builder /api/build/api /api

ENTRYPOINT  ["/api"]
