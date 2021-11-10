package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"sort"
)

type MostMatchFieldHandler struct {
	analysisType int
	countMap     map[string]int
}

func NewMostMatchFieldHandler(analysisType int) *MostMatchFieldHandler {
	return &MostMatchFieldHandler{
		analysisType: analysisType,
		countMap:     make(map[string]int),
	}
}

func (handler *MostMatchFieldHandler) Input(info *ioutil.LogInfo) {
	var field string
	switch handler.analysisType {
	case AnalysisTypeFieldIp:
		field = info.RemoteAddr
	case AnalysisTypeFieldUri:
		field = info.Request
	case AnalysisTypeFieldUserAgent:
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

func (handler *MostMatchFieldHandler) Output(limit int) {
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
