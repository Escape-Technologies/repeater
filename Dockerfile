FROM cgr.dev/chainguard/go:latest as builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

ARG VERSION
ARG COMMIT

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o /bin/prog cmd/repeater/repeater.go

FROM cgr.dev/chainguard/static:latest

USER nonroot

COPY --from=builder /bin/prog /prog

ENTRYPOINT ["/prog"]
