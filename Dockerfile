FROM golang:latest as builder
COPY ./cmd /app/cmd
COPY ./pkg /app/pkg
COPY ./go.* /app/
WORKDIR /app
RUN go build -ldflags="-extldflags=-static" -o mumbleui cmd/mumbleui/main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/mumbleui /app/
COPY ./web /app/web
CMD ["./mumbleui"]