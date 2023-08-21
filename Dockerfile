FROM golang:1.21.0-alpine as builder

WORKDIR /usr/src

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/prog ./cmd/repeater/repeater.go

FROM alpine:3.14
COPY --from=builder /usr/local/bin/prog ./prog
ENTRYPOINT [ "./prog" ]
