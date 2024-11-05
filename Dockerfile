# Stage 1: Modules caching
FROM golang:alpine AS modules
WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Builder
FROM golang:alpine AS builder
WORKDIR /app
COPY --from=modules /go/pkg/ /go/pkg
COPY . .
RUN go build -o unimates ./cmd/main.go
RUN ls -l /app


# Stage 3: Final
FROM alpine
COPY --from=builder /app/unimates .

EXPOSE 8080
CMD ["./unimates"]