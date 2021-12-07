package handler

import (
	"github.com/fantasticmao/nginx-log-analyzer/ioutil"
	"github.com/fantasticmao/nginx-log-analyzer/parser"
	"github.com/pterm/pterm"
	"math"
	"sort"
)

type LargestPercentTimeUrisHandler struct {
	percentile      float64
	timeCostListMap map[string][]float64
}

func NewLargestPercentTimeUrisHandler(percentile float64) *LargestPercentTimeUrisHandler {
	if percentile <= 0 || percentile > 100 {
		ioutil.Fatal("illegal argument percentile: %.3f\n", percentile)
		return nil
	}
	return &LargestPercentTimeUrisHandler{
		percentile:      percentile,
		timeCostListMap: make(map[string][]float64),
	}
}

func (handler *LargestPercentTimeUrisHandler) Input(info *parser.LogInfo) {
	if _, ok := handler.timeCostListMap[info.Request]; ok {
		handler.timeCostListMap[info.Request] = append(handler.timeCostListMap[info.Request], info.RequestTime)
	} else {
		array := []float64{info.RequestTime}
		handler.timeCostListMap[info.Request] = array
	}
}

func (handler *LargestPercentTimeUrisHandler) Output(limit int) {
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

	data := pterm.Bars{}
	for i := 0; i < limit && i < len(keys); i++ {
		data = append(data, pterm.Bar{
			Label: keys[i],
			Value: int(timeCostMap[keys[i]] * 1000),
		})
	}

	ioutil.PTermHeader.Printf("Largest percentile(P%.2f) response time(millisecond) URIs", handler.percentile)
	_ = ioutil.PTermBarChart.WithBars(data).Render()
}
