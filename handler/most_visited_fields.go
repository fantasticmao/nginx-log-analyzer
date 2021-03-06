package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-log-analyzer/ioutil"
	"github.com/fantasticmao/nginx-log-analyzer/parser"
	"sort"
)

type MostVisitedFieldsHandler struct {
	analysisType int
	countMap     map[string]int
}

func NewMostVisitedFieldsHandler(analysisType int) *MostVisitedFieldsHandler {
	return &MostVisitedFieldsHandler{
		analysisType: analysisType,
		countMap:     make(map[string]int),
	}
}

func (handler *MostVisitedFieldsHandler) Input(info *parser.LogInfo) {
	var field string
	switch handler.analysisType {
	case AnalysisTypeVisitedIps:
		field = info.RemoteAddr
	case AnalysisTypeVisitedUris:
		field = info.Request
	case AnalysisTypeVisitedUserAgents:
		field = info.HttpUserAgent
	default:
		ioutil.Fatal("unsupported analysis type: %v\n", handler.analysisType)
		return
	}

	if _, ok := handler.countMap[field]; ok {
		handler.countMap[field]++
	} else {
		handler.countMap[field] = 1
	}
}

func (handler *MostVisitedFieldsHandler) Output(limit int) {
	keys := make([]string, 0, len(handler.countMap))
	for k := range handler.countMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return handler.countMap[keys[i]] > handler.countMap[keys[j]]
	})

	for i := 0; i < limit && i < len(keys); i++ {
		fmt.Printf("\"%v\" hits: %v\n", keys[i], handler.countMap[keys[i]])
	}
}
