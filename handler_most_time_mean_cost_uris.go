package main

import (
	"fmt"
	"sort"
)

type MostTimeMeanCostUrisHandler struct {
	timeCostListMap map[string][]float32
}

func NewMostTimeMeanCostUrisHandler() *MostTimeMeanCostUrisHandler {
	return &MostTimeMeanCostUrisHandler{
		timeCostListMap: make(map[string][]float32),
	}
}

func (handler *MostTimeMeanCostUrisHandler) input(info *LogInfo) {
	if _, ok := handler.timeCostListMap[info.Request]; ok {
		handler.timeCostListMap[info.Request] = append(handler.timeCostListMap[info.Request], info.RequestTime)
	} else {
		array := []float32{info.RequestTime}
		handler.timeCostListMap[info.Request] = array
	}
}

func (handler *MostTimeMeanCostUrisHandler) output(limit int) {
	uriCostMap := make(map[string]float32)
	for uri := range handler.timeCostListMap {
		var sum float32 = 0.0
		var length = 0
		for _, cost := range handler.timeCostListMap[uri] {
			sum += cost
			length++
		}
		uriCostMap[uri] = sum / float32(length)
	}

	keys := make([]string, 0, len(uriCostMap))
	for k := range uriCostMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return uriCostMap[keys[i]] > uriCostMap[keys[j]]
	})

	for i := 0; i < limit && i < len(keys); i++ {
		fmt.Printf("\"%v\" average cost: %v\n", keys[i], uriCostMap[keys[i]])
	}
}
