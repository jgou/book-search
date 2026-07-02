FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY src/go.mod ./
RUN go mod download
COPY src/ ./

# Build the Go application
# CGO_ENABLED=0 ensures a statically linked binary, required for scratch/alpine
# GOOS=linux sets the target OS
RUN CGO_ENABLED=0 GOOS=linux go build -o api main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080

CMD ["./api"]