package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"sort"
)

type MostFrequentStatusHandler struct {
	statusCountMap    map[int]int
	statusUriCountMap map[int]map[string]int
}

func NewMostFrequentStatusHandler() *MostFrequentStatusHandler {
	return &MostFrequentStatusHandler{
		statusCountMap:    make(map[int]int),
		statusUriCountMap: make(map[int]map[string]int),
	}
}

func (handler *MostFrequentStatusHandler) Input(info *ioutil.LogInfo) {
	if _, ok1 := handler.statusCountMap[info.Status]; ok1 {
		handler.statusCountMap[info.Status]++
		if _, ok2 := handler.statusUriCountMap[info.Status][info.Request]; ok2 {
			handler.statusUriCountMap[info.Status][info.Request]++
		} else {
			handler.statusUriCountMap[info.Status][info.Request] = 1
		}
	} else {
		handler.statusCountMap[info.Status] = 1
		handler.statusUriCountMap[info.Status] = make(map[string]int)
		handler.statusUriCountMap[info.Status][info.Request] = 1
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
		fmt.Printf("%v hits: %v\n", status, count)

		uriCountMap := handler.statusUriCountMap[status]
		uriCountKeys := make([]string, 0, len(uriCountMap))
		for k := range uriCountMap {
			uriCountKeys = append(uriCountKeys, k)
		}

		sort.Slice(uriCountKeys, func(i, j int) bool {
			return uriCountMap[uriCountKeys[i]] > uriCountMap[uriCountKeys[j]]
		})
		for i := 0; i < limit && i < len(uriCountKeys); i++ {
			fmt.Printf("  |--\"%v\" hits: %v\n", uriCountKeys[i], uriCountMap[uriCountKeys[i]])
		}
	}
}
