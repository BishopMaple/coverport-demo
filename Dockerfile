FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod .
COPY *.go .
RUN go build -o server .

FROM registry.access.redhat.com/ubi9/ubi-minimal
COPY --from=builder /app/server /usr/local/bin/server
EXPOSE 8080
CMD ["server"]
