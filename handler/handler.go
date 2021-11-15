package handler

import (
	"github.com/fantasticmao/nginx-json-log-analyzer/parser"
)

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
	Input(info *parser.LogInfo)

	Output(limit int)
}
