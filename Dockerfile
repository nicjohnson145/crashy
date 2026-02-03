FROM golang:1.25-alpine AS builder
RUN apk add git ca-certificates
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=builder /src/crashy /bin/crashy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/