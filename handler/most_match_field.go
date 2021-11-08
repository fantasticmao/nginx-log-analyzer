package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"sort"
)

type MostMatchFieldHandler struct {
	analyzeType int
	countMap    map[string]int
}

func NewMostMatchFieldHandler(analyzeType int) *MostMatchFieldHandler {
	return &MostMatchFieldHandler{
		analyzeType: analyzeType,
		countMap:    make(map[string]int),
	}
}

func (handler *MostMatchFieldHandler) Input(info *ioutil.LogInfo) {
	var field string
	switch handler.analyzeType {
	case AnalyzeTypeFieldIp:
		field = info.RemoteAddr
	case AnalyzeTypeFieldUri:
		field = info.Request
	case AnalyzeTypeFieldUserAgent:
		field = info.HttpUserAgent
	default:
		ioutil.Fatal("unsupported analyze type: %v\n", handler.analyzeType)
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
