package parser

import (
	"encoding/json"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
)

const (
	ParseTypeJson = iota
	ParseTypeCombined
)

func JsonLogParser(line []byte) *ioutil.LogInfo {
	logInfo := &ioutil.LogInfo{}
	err := json.Unmarshal(line, logInfo)
	if err != nil {
		ioutil.Fatal("json unmarshal error: %v\n", err.Error())
		return nil
	}
	return logInfo
}

func CombinedLogParser(line []byte) *ioutil.LogInfo {
	var (
		remoteAddr    string
		remoteUser    string
		timeLocal     string
		request       string
		status        int
		bodyBytesSent int
		httpReferer   string
		httpUserAgent string
	)

	return &ioutil.LogInfo{
		RemoteAddr:    remoteAddr,
		RemoteUser:    remoteUser,
		TimeLocal:     timeLocal,
		Request:       request,
		Status:        status,
		BodyBytesSent: bodyBytesSent,
		HttpReferer:   httpReferer,
		HttpUserAgent: httpUserAgent,
	}
}
