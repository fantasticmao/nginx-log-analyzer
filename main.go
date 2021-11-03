package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	showVersion bool
	logFiles    []string
	analyzeType int
	limit       int
	percentile  float64
)

var (
	Name       = "nginx-json-log-analyzer"
	Version    string
	BuildTime  string
	CommitHash string
)

const (
	AnalyzeTypePvUv = iota
	AnalyzeTypeFieldIp
	AnalyzeTypeFieldUri
	AnalyzeTypeFieldUserAgent
	AnalyzeTypeFieldUserCountry
	AnalyzeTypeFieldUserCity
	AnalyzeTypeResponseStatus
	AnalyzeTypeTimeMeanCostUris
	AnalyzeTypeTimePercentCostUris
)

type LogInfo struct {
	TimeIso8601   string  `json:"time_iso8601"`
	RemoteAddr    string  `json:"remote_addr"`
	RequestTime   float64 `json:"request_time"`
	Request       string  `json:"request"`
	Status        int     `json:"status"`
	BodyBytesSent int32   `json:"body_bytes_sent"`
	HttpUserAgent string  `json:"http_user_agent"`
}

type Handler interface {
	input(info *LogInfo)

	output(limit int)
}

func init() {
	flag.BoolVar(&showVersion, "v", false, "show current version")
	flag.IntVar(&analyzeType, "t", 0, "specify the analyze type")
	flag.IntVar(&limit, "n", 15, "limit the number of lines displayed")
	flag.Float64Var(&percentile, "p", 95, "specify the percentile value")
	flag.Parse()
	logFiles = flag.Args()
}

func main() {
	if showVersion {
		fmt.Printf("%v %v build at %v on commit %v\n", Name, Version, BuildTime, CommitHash)
		return
	}

	handler := newHandler(analyzeType)
	process(logFiles, handler)
}

func newHandler(analyzeType int) Handler {
	switch analyzeType {
	case AnalyzeTypePvUv:
		return NewPvAndUvHandler()
	case AnalyzeTypeFieldIp:
		return NewMostMatchFieldHandler(AnalyzeTypeFieldIp)
	case AnalyzeTypeFieldUri:
		return NewMostMatchFieldHandler(AnalyzeTypeFieldUri)
	case AnalyzeTypeFieldUserAgent:
		return NewMostMatchFieldHandler(AnalyzeTypeFieldUserAgent)
	case AnalyzeTypeResponseStatus:
		return NewMostFrequentStatusHandler()
	case AnalyzeTypeTimeMeanCostUris:
		return NewTopTimeMeanCostUrisHandler()
	case AnalyzeTypeTimePercentCostUris:
		return NewTopTimePercentCostUrisHandler(percentile)
	default:
		fatal("unsupported analyze type: %v\n", analyzeType)
		return nil
	}
}

func process(logFiles []string, handler Handler) {
	for _, logFile := range logFiles {
		file, isGzip := openFile(logFile)
		reader := readFile(file, isGzip)
		for {
			data, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fatal("read file error: %v\n", err.Error())
			}

			logInfo := &LogInfo{}
			err = json.Unmarshal(data[:len(data)-1], logInfo)
			if err != nil {
				fatal("json unmarshal error: %v\n", err.Error())
			}
			handler.input(logInfo)
		}
	}

	handler.output(limit)
}

func openFile(path string) (*os.File, bool) {
	file, err := os.Open(path)
	if err != nil {
		fatal("open file error: %v\n", err.Error())
	}

	ext := filepath.Ext(file.Name())
	return file, strings.EqualFold(".gz", ext)
}

func readFile(file *os.File, isGzip bool) *bufio.Reader {
	if isGzip {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			fatal("gzip new reader error: %v\n", err.Error())
		}
		return bufio.NewReader(gzipReader)
	} else {
		return bufio.NewReader(file)
	}
}

func fatal(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a)
	os.Exit(1)
}
