FROM golang:1.18 AS builder
RUN mkdir /app
COPY vendor /app
COPY go.mod /app
COPY go.sum /app
WORKDIR /app
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mainapp ./cmd/main.go

FROM alpine:3.14.2
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/mainapp .

EXPOSE 80 5001 8080 8090
CMD ["./mainapp"]