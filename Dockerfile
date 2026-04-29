FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
COPY internal/ internal/
RUN CGO_ENABLED=0 go build -cover -o server .

FROM registry.access.redhat.com/ubi9/ubi-minimal
RUN mkdir -p /tmp/covdata
ENV GOCOVERDIR=/tmp/covdata
COPY --from=builder /app/server /usr/local/bin/server
EXPOSE 8080
EXPOSE 9095
CMD ["server"]
