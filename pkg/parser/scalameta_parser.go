package parser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/amenzhinsky/go-memexec"
	"github.com/bazelbuild/rules_go/go/tools/bazel"

	sppb "github.com/stackb/scala-gazelle/build/stack/gazelle/scala/parse"
	"github.com/stackb/scala-gazelle/pkg/collections"
)

const (
	contentTypeJSON = "application/json"
	// debugParse is a debug flag for use by a developer
	debugParse = false
)

func NewScalametaParser() *ScalametaParser {
	return &ScalametaParser{}
}

// ScalametaParser is a service that communicates to a scalameta-js parser
// backend over HTTP.
type ScalametaParser struct {
	sppb.UnimplementedParserServer

	process    *memexec.Exec
	processDir string
	cmd        *exec.Cmd

	httpClient *http.Client
	httpUrl    string

	HttpPort int
}

func (s *ScalametaParser) Stop() {
	if s.httpClient != nil {
		s.httpClient.CloseIdleConnections()
		s.httpClient = nil
	}
	if s.cmd != nil {
		s.cmd.Process.Kill()
		s.cmd = nil
	}
	if s.process != nil {
		s.process.Close()
		s.process = nil
	}
	if s.processDir != "" {
		os.RemoveAll(s.processDir)
		s.processDir = ""
	}
}

func (s *ScalametaParser) Start() error {
	t1 := time.Now()

	//
	// Setup temp process directory and write js files
	//
	processDir, err := bazel.NewTmpDir("")
	if err != nil {
		return fmt.Errorf("creating tmp process dir: %w", err)
	}

	scriptPath := filepath.Join(processDir, "scalameta_parser.mjs")
	parserPath := filepath.Join(processDir, "node_modules", "scalameta-parsers", "index.js")

	if err := os.MkdirAll(filepath.Dir(parserPath), os.ModePerm); err != nil {
		return fmt.Errorf("mkdir process tmpdir: %w", err)
	}
	if err := os.WriteFile(scriptPath, []byte(parserrMjs), os.ModePerm); err != nil {
		return fmt.Errorf("writing %s: %w", parserrMjs, err)
	}
	if err := os.WriteFile(parserPath, []byte(scalametaParsersIndexJs), os.ModePerm); err != nil {
		return fmt.Errorf("writing %s: %w", scalametaParsersIndexJs, err)
	}

	if debugParse {
		collections.ListFiles(".")
	}

	//
	// ensure we have a port
	//
	if s.HttpPort == 0 {
		port, err := getFreePort()
		if err != nil {
			return status.Errorf(codes.FailedPrecondition, "getting http port: %v", err)
		}
		s.HttpPort = port
	}
	s.httpUrl = fmt.Sprintf("http://127.0.0.1:%d", s.HttpPort)

	//
	// Setup the bun process
	//
	exe, err := memexec.New(nodeExe)
	if err != nil {
		return err
	}
	s.process = exe

	//
	// Start the bun process
	//
	cmd := exe.Command("scalameta_parser.mjs")
	cmd.Dir = processDir
	cmd.Env = []string{
		"NODE_PATH=" + processDir,
		fmt.Sprintf("PORT=%d", s.HttpPort),
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	s.cmd = cmd

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting process %s: %w", scalametaParsersIndexJs, err)
	}
	go func() {
		// does it make sense to wait for the process?  We kill it forcefully
		// at the end anyway...
		if err := cmd.Wait(); err != nil {
			if err.Error() != "signal: killed" {
				log.Printf("command wait err: %v", err)
			}
		}
	}()

	host := "localhost"
	port := s.HttpPort
	timeout := 3 * time.Second
	if !collections.WaitForConnectionAvailable(host, port, timeout) {
		return fmt.Errorf("waiting to connect to scala parse server %s:%d within %s", host, port, timeout)
	}

	//
	// Setup the http client
	//
	s.httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	t2 := time.Since(t1).Round(1 * time.Millisecond)
	if debugParse {
		log.Printf("parser started (%v)", t2)
	}

	return nil
}

func (s *ScalametaParser) Parse(ctx context.Context, in *sppb.ParseRequest) (*sppb.ParseResponse, error) {
	req, err := newHttpParseRequest(s.httpUrl, in)
	if err != nil {
		return nil, err
	}
	w, err := s.httpClient.Do(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response error: %v", err)
	}

	if debugParse {
		respDump, err := httputil.DumpResponse(w, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("HTTP_RESPONSE:\n%s", string(respDump))
	}

	contentType := w.Header.Get("Content-Type")
	if contentType != contentTypeJSON {
		return nil, status.Errorf(codes.Internal, "response content-type error, want %q, got: %q", contentTypeJSON, contentType)
	}

	data, err := io.ReadAll(w.Body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response data error: %v", err)
	}

	if debugParse {
		log.Printf("response body: %s", string(data))
	}

	var response sppb.ParseResponse
	if err := protojson.Unmarshal(data, &response); err != nil {
		return nil, status.Errorf(codes.Internal, "response body error: %v\n%s", err, string(data))
	}

	return &response, nil
}

// getFreePort asks the kernel for a free open port that is ready to use.
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func newHttpParseRequest(url string, in *sppb.ParseRequest) (*http.Request, error) {
	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "request URL is required")
	}
	if in == nil {
		return nil, status.Errorf(codes.InvalidArgument, "ParseRequest is required")
	}

	json, err := protojson.Marshal(in)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
