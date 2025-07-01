FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY qasm /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/qasm"]