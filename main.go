package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/handler"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"io"
	"os"
	"path"
	"time"
)

var (
	logFiles    []string
	showVersion bool
	configDir   string
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

func init() {
	flag.BoolVar(&showVersion, "v", false, "show current version")
	flag.StringVar(&configDir, "d", "", "specify the configuration directory")
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

	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			ioutil.Fatal("get user home directory error: %v\n", err.Error())
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

	h := handler.NewHandler(configDir, analyzeType, percentile)
	process(logFiles, h, since, util)
}

func process(logFiles []string, h handler.Handler, since, util time.Time) {
nextFile:
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

			// 2. parse json
			logInfo := &ioutil.LogInfo{}
			err = json.Unmarshal(data[:len(data)-1], logInfo)
			if err != nil {
				ioutil.Fatal("json unmarshal error: %v\n", err.Error())
				continue
			}

			// 3. time filter
			logTime, err := time.Parse(time.RFC3339, logInfo.TimeIso8601)
			if err != nil {
				ioutil.Fatal("parse log time error: %v\n", err.Error())
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
			h.Input(logInfo)
		}

		// 5. close file handler
		err := file.Close()
		if err != nil {
			ioutil.Fatal("close file error: %v\n", err.Error())
		}
	}

	// 5. print result
	h.Output(limit)
}
