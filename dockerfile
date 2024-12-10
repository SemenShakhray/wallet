FROM golang:1.22.1-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY cmd ./cmd
COPY internal ./internal
COPY tests ./tests
RUN go build -o ./app ./cmd/main.go

FROM alpine AS runner

COPY --from=builder /app /
COPY /config.env /config.env


CMD ["/app"]