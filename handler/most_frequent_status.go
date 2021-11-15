package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/parser"
	"sort"
)

type MostFrequentStatusHandler struct {
	// status -> count
	statusCountMap map[int]int
	// status -> uri -> count
	statusUriCountMap map[int]map[string]int
}

func NewMostFrequentStatusHandler() *MostFrequentStatusHandler {
	return &MostFrequentStatusHandler{
		statusCountMap:    make(map[int]int),
		statusUriCountMap: make(map[int]map[string]int),
	}
}

func (handler *MostFrequentStatusHandler) Input(info *parser.LogInfo) {
	if _, ok := handler.statusUriCountMap[info.Status]; !ok {
		handler.statusCountMap[info.Status] = 1
		handler.statusUriCountMap[info.Status] = make(map[string]int)
	} else {
		handler.statusCountMap[info.Status]++
	}

	if _, ok := handler.statusUriCountMap[info.Status][info.Request]; !ok {
		handler.statusUriCountMap[info.Status][info.Request] = 1
	} else {
		handler.statusUriCountMap[info.Status][info.Request]++
	}
}

func (handler *MostFrequentStatusHandler) Output(limit int) {
	statusCountKeys := make([]int, 0, len(handler.statusCountMap))
	for k := range handler.statusCountMap {
		statusCountKeys = append(statusCountKeys, k)
	}
	sort.Ints(statusCountKeys)

	for _, status := range statusCountKeys {
		count := handler.statusCountMap[status]
		uriCountMap := handler.statusUriCountMap[status]
		fmt.Printf("%v hits: %v\n", status, count)

		uriCountKeys := make([]string, 0, len(uriCountMap))
		for k := range uriCountMap {
			uriCountKeys = append(uriCountKeys, k)
		}
		sort.Slice(uriCountKeys, func(i, j int) bool {
			return uriCountMap[uriCountKeys[i]] > uriCountMap[uriCountKeys[j]]
		})

		for i := 0; i < limit && i < len(uriCountKeys); i++ {
			uri := uriCountKeys[i]
			fmt.Printf("  |--\"%v\" hits: %v\n", uri, uriCountMap[uri])
		}
	}
}
