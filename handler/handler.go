package handler

import "github.com/fantasticmao/nginx-json-log-analyzer/ioutil"

const (
	AnalyzeTypePvUv = iota
	AnalyzeTypeFieldIp
	AnalyzeTypeFieldUri
	AnalyzeTypeFieldUserAgent
	AnalyzeTypeFieldUserCity
	AnalyzeTypeResponseStatus
	AnalyzeTypeTimeMeanCostUris
	AnalyzeTypeTimePercentCostUris
)

type Handler interface {
	Input(info *ioutil.LogInfo)

	Output(limit int)
}
