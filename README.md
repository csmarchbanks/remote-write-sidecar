# Remote Write Sidecar

The remote write sidecar implements the remote write code from
[Prometheus](https://prometheus.io) in a separate process from Prometheus. By
removing remote write from the main Prometheus process it is possible scale the
resources of each process separately, and if an error occurs in the remote write
code it will be isolated from the Prometheus process.

## Configuration

### Flags

The remote write sidecar uses a subset of the Prometheus flags:
```
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
      --version                  Show application version.
      --config.file="prometheus.yml"
                                 Prometheus configuration file path.
      --web.listen-address="0.0.0.0:9095"
                                 Address to listen on for UI, API, and telemetry.
      --web.read-timeout=5m      Maximum duration before timing out read of the request, and closing idle connections.
      --web.max-connections=512  Maximum number of simultaneous connections.
      --storage.tsdb.path="data/"
                                 Base path for metrics storage.
      --storage.remote.flush-deadline=<duration>
                                 How long to wait flushing sample on shutdown or config reload.
      --log.level=info           Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt        Output format of log messages. One of: [logfmt, json]
```

### Configuration file

The remote write sidecar can be configured using the same YAML configuration
file as Prometheus. The important stanzas are:
```
global:
  external_labels:
    foo: bar
remote_write
  - url: https://remote-write-target
```
