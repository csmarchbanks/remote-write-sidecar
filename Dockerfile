# Build container
FROM golang:1.12.5
WORKDIR /go/src/github.com/csmarchbanks/remote-write-sidecar
COPY . .
RUN CGO_ENABLED=0 make build

FROM quay.io/prometheus/busybox:latest
EXPOSE 9095
USER nobody
COPY --from=0 /go/src/github.com/csmarchbanks/remote-write-sidecar/remotewrite /bin/remotewrite
ENTRYPOINT ["/bin/remotewrite"]
CMD        [ "--config.file=/etc/remotewrite/remotewrite.yml", \
             "--storage.tsdb.path=/prometheus" ]
