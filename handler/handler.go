package handler

import "github.com/fantasticmao/nginx-json-log-analyzer/ioutil"

const (
	AnalysisTypePvUv = iota
	AnalysisTypeFieldIp
	AnalysisTypeFieldUri
	AnalysisTypeFieldUserAgent
	AnalysisTypeFieldUserCity
	AnalysisTypeResponseStatus
	AnalysisTypeTimeMeanCostUris
	AnalysisTypeTimePercentCostUris
)

type Handler interface {
	Input(info *ioutil.LogInfo)

	Output(limit int)
}
