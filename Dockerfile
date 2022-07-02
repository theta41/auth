FROM golang:1.18-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o application cmd/auth/main.go

FROM alpine:3.15.4
ENV DB_PASSWORD ${DB_PASSWORD}
ENV DB_LOGIN ${DB_LOGIN}
WORKDIR /app
COPY --from=builder /app/application /app/application
COPY *.yaml ./
COPY *.json ./
CMD ["/app/application", "-c", "/app/config.yaml"]