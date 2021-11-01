package main

import (
	"fmt"
	"sort"
)

type MostTimeCostUrisHandler struct {
	uriCostListMap map[string][]float32
}

func NewMostTimeCostUrisHandler() *MostTimeCostUrisHandler {
	return &MostTimeCostUrisHandler{
		uriCostListMap: make(map[string][]float32),
	}
}

func (handler *MostTimeCostUrisHandler) input(info *LogInfo) {
	if _, ok := handler.uriCostListMap[info.Request]; ok {
		handler.uriCostListMap[info.Request] = append(handler.uriCostListMap[info.Request], info.RequestTime)
	} else {
		array := []float32{info.RequestTime}
		handler.uriCostListMap[info.Request] = array
	}
}

func (handler *MostTimeCostUrisHandler) output(limit int) {
	uriCostMap := make(map[string]float32)
	for uri := range handler.uriCostListMap {
		var sum float32 = 0.0
		var length = 0
		for _, cost := range handler.uriCostListMap[uri] {
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
