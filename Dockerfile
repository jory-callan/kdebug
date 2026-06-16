# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -ldflags="-s -w" -o kdebug .

# Runtime stage
FROM alpine:3.19
ENV TZ=Asia/Shanghai
WORKDIR /app
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
    && apk add --no-cache ca-certificates
COPY --from=builder /app/kdebug /app/kdebug
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/kdebug"]
CMD ["-c", "echo"]