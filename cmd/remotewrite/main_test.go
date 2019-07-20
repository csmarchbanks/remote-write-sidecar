package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/tsdb"
	"github.com/prometheus/tsdb/labels"
	"github.com/prometheus/tsdb/testutil"
)

type sample struct {
	Labels    labels.Labels
	Timestamp int64
	Value     float64
}

var (
	path string
)

func TestMain(m *testing.M) {
	flag.Parse()

	var err error
	path, err = os.Getwd()
	if err != nil {
		fmt.Printf("can't get current dir :%s \n", err)
		os.Exit(1)
	}
	path = filepath.Join(path, "remotewrite")

	build := exec.Command("go", "build", "-o", path)
	output, err := build.CombinedOutput()
	if err != nil {
		fmt.Printf("compilation error :%s \n", output)
		os.Exit(1)
	}

	exitCode := m.Run()
	os.Remove(path)
	os.Exit(exitCode)
}

func TestStartupInterrupt(t *testing.T) {
	configPath := filepath.Join("..", "..", "examples", "remotewrite.yml")
	remotewrite := exec.Command(path, "--config.file="+configPath)
	testutil.Ok(t, remotewrite.Start())

	done := make(chan error)
	go func() {
		done <- remotewrite.Wait()
	}()

	var stoppedErr error

	if !isUp() {
		t.Fatalf("remotewrite didn't start in the specified timeout")
	}

	remotewrite.Process.Signal(os.Interrupt)
	select {
	case stoppedErr = <-done:
	case <-time.After(10 * time.Second):
	}

	if err := remotewrite.Process.Kill(); err == nil {
		t.Errorf("remotewrite didn't shutdown gracefully after sending an interrupt")
	} else if stoppedErr != nil && stoppedErr.Error() != "signal: interrupt" {
		t.Errorf("remotewrite exited with an unexpected error:%v", stoppedErr)
	}
}

func TestWriteSamples(t *testing.T) {
	nSamples := 1000
	sampleCh := make(chan sample, nSamples)
	remoteStorage := remoteStorageServer(sampleCh)
	defer remoteStorage.Close()

	configYaml := fmt.Sprintf(`
global:
  external_labels:
    prometheus: test
remote_write:
- url: %s/receive
  queue_config:
    batch_send_deadline: 500ms
`, remoteStorage.URL)
	ioutil.WriteFile("testconfig.yml", []byte(configYaml), 0644)
	defer os.Remove("testconfig.yml")

	dataPath := "data/"
	db, err := tsdb.Open(dataPath, nil, nil, tsdb.DefaultOptions)
	testutil.Ok(t, err)
	defer os.RemoveAll(dataPath)
	defer db.Close()

	remotewrite := exec.Command(path,
		"--config.file=testconfig.yml",
		"--storage.tsdb.path="+dataPath,
	)
	var errb bytes.Buffer
	remotewrite.Stderr = &errb
	err = remotewrite.Start()
	testutil.Ok(t, err)
	defer remotewrite.Process.Kill()

	if !isUp() {
		t.Fatalf("remotewrite didn't start in the specified timeout")
	}

	lbls := labels.FromStrings("hello", "world")
	app := db.Appender()
	for i := 0; i < nSamples; i++ {
		app.Add(lbls, time.Now().UnixNano()/1e6, 1.0)
	}
	testutil.Ok(t, app.Commit())

	expectedLabels := labels.FromStrings("hello", "world", "prometheus", "test")
	timeout := time.After(60 * time.Second)
	for i := 0; i < nSamples; i++ {
		select {
		case sample := <-sampleCh:
			testutil.Equals(t, expectedLabels, sample.Labels)
		case <-timeout:
			t.Errorf("timeout waiting for samples reached, received %d / %d samples", i, nSamples)
		}
	}
}

func remoteStorageServer(samples chan sample) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/receive",
		func(w http.ResponseWriter, r *http.Request) {
			compressed, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			reqBuf, err := snappy.Decode(nil, compressed)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var req prompb.WriteRequest
			if err := proto.Unmarshal(reqBuf, &req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			for _, ts := range req.Timeseries {
				lbls := make(labels.Labels, 0, len(ts.Labels))
				for _, l := range ts.Labels {
					lbls = append(lbls, labels.Label{Name: l.Name, Value: l.Value})
				}

				for _, s := range ts.Samples {
					samples <- sample{
						Labels:    lbls,
						Timestamp: s.Timestamp,
						Value:     s.Value,
					}
				}
			}
		},
	)

	return httptest.NewServer(mux)
}

func isUp() bool {
	for x := 0; x < 10; x++ {
		if _, err := http.Get("http://localhost:9095/metrics"); err == nil {
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}
