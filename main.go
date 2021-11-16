package main

import (
	"flag"
	"fmt"
	"github.com/fantasticmao/nginx-log-analyzer/handler"
	"github.com/fantasticmao/nginx-log-analyzer/ioutil"
	"github.com/fantasticmao/nginx-log-analyzer/parser"
	"io"
	"os"
	"path"
	"time"
)

var (
	logFiles     []string
	showVersion  bool
	configDir    string
	analysisType int
	limit        int
	limitSecond  int
	percentile   float64
	timeAfter    string
	timeBefore   string
	logFormat    string
)

var (
	Name       = "nginx-log-analyzer"
	Version    string
	BuildTime  string
	CommitHash string
)

func init() {
	flag.BoolVar(&showVersion, "v", false, "show current version")
	flag.StringVar(&configDir, "d", "", "specify the configuration directory")
	flag.IntVar(&analysisType, "t", 0, "specify the analysis type, see documentation for more details:\nhttps://github.com/fantasticmao/nginx-log-analyzer#specify-the-analysis-type--t")
	flag.IntVar(&limit, "n", 15, "limit the output lines number")
	flag.IntVar(&limitSecond, "n2", 15, "limit the secondary output lines number in '-t 4' mode")
	flag.Float64Var(&percentile, "p", 95, "specify the percentile value in '-t 7' mode")
	flag.StringVar(&timeAfter, "ta", "", "limit the analysis start time, in format of RFC3339 e.g. '2021-11-01T00:00:00+08:00'")
	flag.StringVar(&timeBefore, "tb", "", "limit the analysis end time, in format of RFC3339 e.g. '2021-11-02T00:00:00+08:00'")
	flag.StringVar(&logFormat, "lf", "combined", "specify the nginx log format, value should be 'combined' or 'json'")
	flag.Parse()
	logFiles = flag.Args()
}

func main() {
	if showVersion {
		fmt.Printf("%v %v build at %v on commit %v\n", Name, Version, BuildTime, CommitHash)
		return
	}

	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			ioutil.Fatal("get user home directory error: %v\n", err.Error())
			return
		}
		configDir = path.Join(homeDir, ".config", Name)
	}

	var (
		since, util time.Time
		err         error
	)
	if timeAfter != "" {
		since, err = time.Parse(time.RFC3339, timeAfter)
		if err != nil {
			ioutil.Fatal("parse start time error: %v\n", err.Error())
			return
		}
	}
	if timeBefore != "" {
		util, err = time.Parse(time.RFC3339, timeBefore)
		if err != nil {
			ioutil.Fatal("parse end time error: %v\n", err.Error())
			return
		}
	}

	p := newLogParser()
	h := newLogHandler()
	process(logFiles, p, h, since, util)
}

func newLogHandler() handler.Handler {
	switch analysisType {
	case handler.AnalysisTypePvAndUv:
		return handler.NewPvAndUvHandler()
	case handler.AnalysisTypeVisitedIps:
		return handler.NewMostVisitedFieldsHandler(analysisType)
	case handler.AnalysisTypeVisitedUris:
		return handler.NewMostVisitedFieldsHandler(analysisType)
	case handler.AnalysisTypeVisitedUserAgents:
		return handler.NewMostVisitedFieldsHandler(analysisType)
	case handler.AnalysisTypeVisitedLocations:
		const dbFile = "City.mmdb"
		return handler.NewMostVisitedLocationsHandler(path.Join(configDir, dbFile), limitSecond)
	case handler.AnalysisTypeResponseStatus:
		return handler.NewMostFrequentStatusHandler()
	case handler.AnalysisTypeAverageTimeUris:
		return handler.NewLargestAverageTimeUrisHandler()
	case handler.AnalysisTypePercentTimeUris:
		return handler.NewLargestPercentTimeUrisHandler(percentile)
	default:
		ioutil.Fatal("unsupported analysis type: %v\n", analysisType)
		return nil
	}
}

func newLogParser() parser.Parser {
	switch logFormat {
	case parser.LogFormatTypeCombined:
		return parser.NewCombinedParser()
	case parser.LogFormatTypeJson:
		return parser.NewJsonParser()
	default:
		ioutil.Fatal("unsupported log format : %v\n", logFormat)
		return nil
	}
}

func process(logFiles []string, p parser.Parser, h handler.Handler, since, util time.Time) {
	for _, logFile := range logFiles {
		// 1. open and read file
		file, isGzip := ioutil.OpenFile(logFile)
		reader := ioutil.ReadFile(file, isGzip)
		for {
			data, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				ioutil.Fatal("read file error: %v\n", err.Error())
				return
			}

			// 2. parse line
			logInfo := p.ParseLog(data)

			// 3. datetime filter
			if !since.IsZero() || !util.IsZero() {
				logTime := parser.ParseTime(logInfo.TimeLocal)
				if !since.IsZero() && logTime.Before(since) {
					// go to next line
					continue
				}
				if !util.IsZero() && logTime.After(util) {
					// go to next file
					break
				}
			}

			// 4. process data
			h.Input(logInfo)
		}

		// 5. close file handler
		err := file.Close()
		if err != nil {
			ioutil.Fatal("close file error: %v\n", err.Error())
			return
		}
	}

	// 5. print result
	h.Output(limit)
}
