FROM golang:alpine AS builder
WORKDIR /backend
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go get -u github.com/pressly/goose/v3/cmd/goose
RUN go build -o /go/bin/goose github.com/pressly/goose/v3/cmd/goose

RUN go build -o unimates ./cmd/main.go


FROM alpine:latest
WORKDIR /backend

COPY --from=builder /backend/unimates .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

EXPOSE 8080
CMD ["./unimates"]