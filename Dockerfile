FROM golang:1.21.0-alpine as builder

WORKDIR /usr/src

COPY go.mod go.sum ./

RUN go mod download && go mod verify

ARG VERSION
ARG COMMIT

COPY . .
RUN go build -ldflags="-s -w -X main.version=$VERSION -X main.commit=$COMMIT" -v -o /usr/local/bin/prog ./cmd/repeater/repeater.go

FROM alpine:3.14

RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/prog ./prog

ENTRYPOINT [ "sh", "-c", "update-ca-certificates && ./prog" ]
