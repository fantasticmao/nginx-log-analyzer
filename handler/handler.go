package handler

import (
	"github.com/fantasticmao/nginx-log-analyzer/parser"
)

const (
	AnalysisTypePvAndUv = iota
	AnalysisTypeVisitedIps
	AnalysisTypeVisitedUris
	AnalysisTypeVisitedUserAgents
	AnalysisTypeVisitedLocations
	AnalysisTypeResponseStatus
	AnalysisTypeAverageTimeUris
	AnalysisTypePercentTimeUris
)

type Handler interface {
	Input(info *parser.LogInfo)

	Output(limit int)
}
