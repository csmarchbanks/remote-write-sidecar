module github.com/csmarchbanks/remote-write-sidecar

go 1.12

require (
	github.com/aws/aws-sdk-go v1.19.41 // indirect
	github.com/go-kit/kit v0.8.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/snappy v0.0.1
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/miekg/dns v1.1.13 // indirect
	github.com/oklog/run v1.0.0
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.4.1
	github.com/prometheus/prometheus v2.11.0+incompatible
	github.com/prometheus/tsdb v0.9.1
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec // indirect
	golang.org/x/net v0.0.0-20190403144856-b630fd6fe46b
	golang.org/x/oauth2 v0.0.0-20190523182746-aaccbc9213b0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/api v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20190530194941-fb225487d101 // indirect
	google.golang.org/grpc v1.21.0 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7 // indirect
)

replace k8s.io/klog => github.com/simonpasquier/klog-gokit v0.1.0
