FROM alpine:latest

COPY containerd /app

EXPOSE 8080

ENTRYPOINT ["./app"]