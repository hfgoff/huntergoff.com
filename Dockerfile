ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod tidy && go mod download && go mod verify

COPY . .
ENV CGO_ENABLED=0
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build"  go build -v -o /run-app .

FROM debian:bookworm

# To make HTTP(s) requetss work 
RUN apt-get update && apt-get install -y ca-certificates &&  rm -rf /var/lib/apt/lists/*

COPY --from=builder /run-app /usr/local/bin/
COPY templates/ /templates/
CMD ["run-app"]
