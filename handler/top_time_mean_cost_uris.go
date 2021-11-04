package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"sort"
)

type TopTimeMeanCostUrisHandler struct {
	timeCostListMap map[string][]float64
}

func NewTopTimeMeanCostUrisHandler() *TopTimeMeanCostUrisHandler {
	return &TopTimeMeanCostUrisHandler{
		timeCostListMap: make(map[string][]float64),
	}
}

func (handler *TopTimeMeanCostUrisHandler) Input(info *ioutil.LogInfo) {
	if _, ok := handler.timeCostListMap[info.Request]; ok {
		handler.timeCostListMap[info.Request] = append(handler.timeCostListMap[info.Request], info.RequestTime)
	} else {
		array := []float64{info.RequestTime}
		handler.timeCostListMap[info.Request] = array
	}
}

func (handler *TopTimeMeanCostUrisHandler) Output(limit int) {
	timeCostMap := make(map[string]float64)
	for uri, costList := range handler.timeCostListMap {
		var sum = 0.0
		for _, cost := range costList {
			sum += cost
		}
		timeCostMap[uri] = sum / float64(len(costList))
	}

	keys := make([]string, 0, len(timeCostMap))
	for k := range timeCostMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return timeCostMap[keys[i]] > timeCostMap[keys[j]]
	})

	for i := 0; i < limit && i < len(keys); i++ {
		fmt.Printf("\"%v\" mean response-time: %.3f\n", keys[i], timeCostMap[keys[i]])
	}
}
