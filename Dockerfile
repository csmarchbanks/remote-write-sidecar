# Build container
FROM golang:1.12.5
WORKDIR /go/src/github.com/csmarchbanks/remote-write-sidecar
COPY . .
RUN CGO_ENABLED=0 make build

FROM alpine:3.9.2
RUN apk update && apk add --no-cache \
    ca-certificates
EXPOSE 9095
COPY --from=0 /go/src/github.com/csmarchbanks/remote-write-sidecar/remotewrite /bin/remotewrite
ENTRYPOINT ["/bin/remotewrite"]
CMD        [ "--config.file=/etc/remotewrite/remotewrite.yml", \
             "--storage.tsdb.path=/prometheus" ]

