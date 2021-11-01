package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	showVersion bool
	logFile     string
	limit       int
	analyzeType int
)

var (
	Version    string
	BuildTime  string
	CommitHash string
)

const (
	AnalyzeTypePvAndUv = iota
	AnalyzeTypeFieldIp
	AnalyzeTypeFieldUri
	AnalyzeTypeFieldUserAgent
	AnalyzeTypeTimeCostUris
)

type LogInfo struct {
	TimeIso8601   string  `json:"time_iso8601"`
	RemoteAddr    string  `json:"remote_addr"`
	RequestTime   float32 `json:"request_time"`
	Request       string  `json:"request"`
	Status        int     `json:"status"`
	BodyBytesSent int32   `json:"body_bytes_sent"`
	HttpUserAgent string  `json:"http_user_agent"`
}

type Handler interface {
	input(*LogInfo)

	output(int)
}

func init() {
	flag.BoolVar(&showVersion, "v", false, "show current version")
	flag.StringVar(&logFile, "f", "", "specify JSON log file")
	flag.IntVar(&limit, "n", 20, "limit lines displayed")
	flag.IntVar(&analyzeType, "t", 0, "analyze type")
	flag.Parse()
}

func main() {
	if showVersion {
		fmt.Printf("nginx-json-log-analyze %v build at %v on commit %v\n", Version, BuildTime, CommitHash)
		return
	}

	handler := newHandler(analyzeType)
	process(logFile, handler)
}

func newHandler(analyzeType int) Handler {
	switch analyzeType {
	case AnalyzeTypePvAndUv:
		return NewPvAndUvHandler()
	case AnalyzeTypeFieldIp:
		return NewMostMatchFieldHandler(AnalyzeTypeFieldIp)
	case AnalyzeTypeFieldUri:
		return NewMostMatchFieldHandler(AnalyzeTypeFieldUri)
	case AnalyzeTypeFieldUserAgent:
		return NewMostMatchFieldHandler(AnalyzeTypeFieldUserAgent)
	case AnalyzeTypeTimeCostUris:
		return NewMostTimeCostUrisHandler()
	default:
		panic(errors.New("unknown analyze type"))
	}
}

func process(path string, handler Handler) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var data []byte
	var logInfo *LogInfo
	reader := bufio.NewReader(file)
	for {
		data, err = reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		logInfo = &LogInfo{}
		err = json.Unmarshal(data[:len(data)-1], logInfo)
		if err != nil {
			panic(err)
		}
		handler.input(logInfo)
	}
	handler.output(limit)
}
