package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"math"
	"sort"
)

type TopTimePercentCostUrisHandler struct {
	percentile      float64
	timeCostListMap map[string][]float64
}

func NewTopTimePercentCostUrisHandler(percentile float64) *TopTimePercentCostUrisHandler {
	if percentile <= 0 || percentile > 100 {
		ioutil.Fatal("illegal argument percentile: %.3f\n", percentile)
		return nil
	}
	return &TopTimePercentCostUrisHandler{
		percentile:      percentile,
		timeCostListMap: make(map[string][]float64),
	}
}

func (handler *TopTimePercentCostUrisHandler) Input(info *ioutil.LogInfo) {
	if _, ok := handler.timeCostListMap[info.Request]; ok {
		handler.timeCostListMap[info.Request] = append(handler.timeCostListMap[info.Request], info.RequestTime)
	} else {
		array := []float64{info.RequestTime}
		handler.timeCostListMap[info.Request] = array
	}
}

func (handler *TopTimePercentCostUrisHandler) Output(limit int) {
	timeCostMap := make(map[string]float64)
	for uri, costList := range handler.timeCostListMap {
		sort.Float64s(costList)

		// according to https://stackoverflow.com/questions/41413544/calculate-percentile-from-a-long-array
		index := int(math.Ceil(handler.percentile/100*float64(len(costList))) - 1)
		timeCostMap[uri] = costList[index]
	}

	keys := make([]string, 0, len(timeCostMap))
	for k := range timeCostMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return timeCostMap[keys[i]] > timeCostMap[keys[j]]
	})

	for i := 0; i < limit && i < len(keys); i++ {
		fmt.Printf("\"%v\" P%.2f response-time: %.3f\n", keys[i], handler.percentile, timeCostMap[keys[i]])
	}
}
