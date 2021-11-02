package main

import (
	"fmt"
	"sort"
)

type TopTimeMeanCostUrisHandler struct {
	timeCostListMap map[string][]float32
}

func NewTopTimeMeanCostUrisHandler() *TopTimeMeanCostUrisHandler {
	return &TopTimeMeanCostUrisHandler{
		timeCostListMap: make(map[string][]float32),
	}
}

func (handler *TopTimeMeanCostUrisHandler) input(info *LogInfo) {
	if _, ok := handler.timeCostListMap[info.Request]; ok {
		handler.timeCostListMap[info.Request] = append(handler.timeCostListMap[info.Request], info.RequestTime)
	} else {
		array := []float32{info.RequestTime}
		handler.timeCostListMap[info.Request] = array
	}
}

func (handler *TopTimeMeanCostUrisHandler) output(limit int) {
	timeCostMap := make(map[string]float32)
	for uri, costList := range handler.timeCostListMap {
		var sum float32 = 0.0
		var length = 0
		for _, cost := range costList {
			sum += cost
			length++
		}
		timeCostMap[uri] = sum / float32(length)
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
