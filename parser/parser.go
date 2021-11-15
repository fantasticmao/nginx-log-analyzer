package parser

import (
	"bytes"
	"encoding/json"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"strconv"
)

const (
	LogFormatCombined = "combined"
	LogFormatJson     = "json"
)

type Parser interface {
	Parse(line []byte) *ioutil.LogInfo
}

type JsonParser struct {
}

func NewJsonParser() *JsonParser {
	return &JsonParser{}
}

func (parser *JsonParser) Parse(line []byte) *ioutil.LogInfo {
	logInfo := &ioutil.LogInfo{}
	err := json.Unmarshal(line[:len(line)-1], logInfo)
	if err != nil {
		ioutil.Fatal("parse json log error: %v\n", err.Error())
		return nil
	}
	return logInfo
}

type CombinedParser struct {
	delimiters [][]byte
}

func NewCombinedParser() *CombinedParser {
	//                            0             1             2           3         4           5               6                 7
	// log_format combined '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"';
	//                                   |             |             |          |       |                |               |                  |
	//                                 ' - '          ' ['         '] "'      '" '     ' '              ' "'           '" "'              '"\n'
	var delimiters = [][]byte{
		[]byte(" - "), []byte(" ["), []byte("] \""), []byte("\" "),
		[]byte(" "), []byte(" \""), []byte("\" \""), []byte("\"\n"),
	}
	return &CombinedParser{
		delimiters: delimiters,
	}
}

func (parser *CombinedParser) Parse(line []byte) *ioutil.LogInfo {
	var (
		variables = make([]string, 0, 8)
		i         = 0 // variable start index
		j         = 0 // variable end index
		k         = 0 // delimiters and variables index
	)
	for j < len(line) && k < len(parser.delimiters) {
		if bytes.Equal(line[j:j+len(parser.delimiters[k])], parser.delimiters[k]) {
			variables = append(variables, string(line[i:j]))
			j = j + len(parser.delimiters[k])
			i = j
			k++
		} else {
			j++
		}
	}
	if k != len(parser.delimiters) {
		ioutil.Fatal("parse combined log error: %v\n", string(line))
	}
	status, err := strconv.Atoi(variables[4])
	if err != nil {
		ioutil.Fatal("convert $status to int error: %v\n", variables[4])
	}
	bodyBytesSent, err := strconv.Atoi(variables[5])
	if err != nil {
		ioutil.Fatal("convert $body_bytes_sent to int error: %v\n", variables[5])
	}
	return &ioutil.LogInfo{
		RemoteAddr:    variables[0],
		RemoteUser:    variables[1],
		TimeLocal:     variables[2],
		Request:       variables[3],
		Status:        status,
		BodyBytesSent: bodyBytesSent,
		HttpReferer:   variables[6],
		HttpUserAgent: variables[7],
	}
}
