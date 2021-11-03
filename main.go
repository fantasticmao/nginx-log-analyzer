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
	"time"
)

var (
	showVersion bool
	logFiles    []string
	analyzeType int
	limit       int
	percentile  float64
	timeAfter   string
	timeBefore  string
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
	flag.IntVar(&limit, "n", 15, "limit the output number lines")
	flag.Float64Var(&percentile, "p", 95, "specify the percentile value in '-t 8' mode")
	flag.StringVar(&timeAfter, "ta", "", "specify the analyze start time, in format of '2006-01-02T15:04:05Z07:00'")
	flag.StringVar(&timeBefore, "tb", "", "specify the analyze end time, in format of '2006-01-02T15:04:05Z07:00'")
	flag.Parse()
	logFiles = flag.Args()
}

func main() {
	if showVersion {
		fmt.Printf("%v %v build at %v on commit %v\n", Name, Version, BuildTime, CommitHash)
		return
	}

	var (
		since, util time.Time
		err         error
	)
	if timeAfter != "" {
		since, err = time.Parse(time.RFC3339, timeAfter)
		if err != nil {
			fatal("parse start time error: %v\n", err.Error())
			return
		}
	}
	if timeBefore != "" {
		util, err = time.Parse(time.RFC3339, timeBefore)
		if err != nil {
			fatal("parse end time error: %v\n", err.Error())
			return
		}
	}

	handler := newHandler(analyzeType)
	process(logFiles, handler, since, util)
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

func process(logFiles []string, handler Handler, since, util time.Time) {
nextFile:
	for _, logFile := range logFiles {
		// 1. open and read file
		file, isGzip := openFile(logFile)
		reader := readFile(file, isGzip)
		for {
			data, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fatal("read file error: %v\n", err.Error())
				return
			}

			// 2. parse json
			logInfo := &LogInfo{}
			err = json.Unmarshal(data[:len(data)-1], logInfo)
			if err != nil {
				fatal("json unmarshal error: %v\n", err.Error())
				continue
			}

			// 3. time filter
			logTime, err := time.Parse(time.RFC3339, logInfo.TimeIso8601)
			if err != nil {
				fatal("parse log time error: %v\n", err.Error())
				continue
			}
			if !since.IsZero() && logTime.Before(since) {
				// go to next line
				continue
			}
			if !util.IsZero() && logTime.After(util) {
				// go to next file
				break nextFile
			}

			// 4. process data
			handler.input(logInfo)
		}
	}

	// 5. print result
	handler.output(limit)
}

func openFile(path string) (*os.File, bool) {
	file, err := os.Open(path)
	if err != nil {
		fatal("open file error: %v\n", err.Error())
		return nil, false
	}

	ext := filepath.Ext(file.Name())
	return file, strings.EqualFold(".gz", ext)
}

func readFile(file *os.File, isGzip bool) *bufio.Reader {
	if isGzip {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			fatal("gzip new reader error: %v\n", err.Error())
			return nil
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
