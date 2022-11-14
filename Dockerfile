FROM golang:1.19-alpine3.16 as builder
WORKDIR /go/src/app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  \
    go build -ldflags '-w -s -extldflags "-static"' -o /app ./cmd/bot

FROM scratch
COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app"]
