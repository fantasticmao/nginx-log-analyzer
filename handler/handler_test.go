package handler

import (
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	limit = 15

	ip1 = "192.168.1.1"
	ip2 = "192.168.1.2"
	ip3 = "192.168.1.3"

	uri1 = "GET /name/Tom HTTP/2.0"
	uri2 = "GET /name/Sam HTTP/2.0"
	uri3 = "GET /name/Bob HTTP/2.0"

	userAgent1 = "iOS"
	userAgent2 = "Android"
	userAgent3 = "Windows"

	responseStatus1 = 200
	responseStatus2 = 302
	responseStatus3 = 404

	responseTime1 = 0.1
	responseTime2 = 0.2
	responseTime3 = 0.3
)

func TestPvAndUv(t *testing.T) {
	handler := NewPvAndUvHandler()
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip1})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip1})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip1})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip2})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip2})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip3})
	handler.Output(limit)
	assert.Equal(t, int32(6), handler.pv)
	assert.Equal(t, int32(3), handler.uv)
}

func TestMostMatchFieldIp(t *testing.T) {
	handler := NewMostMatchFieldHandler(AnalyzeTypeFieldIp)
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip1})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip1})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip1})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip2})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip2})
	handler.Input(&ioutil.LogInfo{RemoteAddr: ip3})
	handler.Output(limit)
	assert.Equal(t, 3, handler.countMap[ip1])
	assert.Equal(t, 2, handler.countMap[ip2])
	assert.Equal(t, 1, handler.countMap[ip3])
}

func TestMostMatchFieldUri(t *testing.T) {
	handler := NewMostMatchFieldHandler(AnalyzeTypeFieldUri)
	handler.Input(&ioutil.LogInfo{Request: uri1})
	handler.Input(&ioutil.LogInfo{Request: uri1})
	handler.Input(&ioutil.LogInfo{Request: uri1})
	handler.Input(&ioutil.LogInfo{Request: uri2})
	handler.Input(&ioutil.LogInfo{Request: uri2})
	handler.Input(&ioutil.LogInfo{Request: uri3})
	handler.Output(limit)
	assert.Equal(t, 3, handler.countMap[uri1])
	assert.Equal(t, 2, handler.countMap[uri2])
	assert.Equal(t, 1, handler.countMap[uri3])
}

func TestMostMatchFieldUserAgent(t *testing.T) {
	handler := NewMostMatchFieldHandler(AnalyzeTypeFieldUserAgent)
	handler.Input(&ioutil.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&ioutil.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&ioutil.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&ioutil.LogInfo{HttpUserAgent: userAgent2})
	handler.Input(&ioutil.LogInfo{HttpUserAgent: userAgent2})
	handler.Input(&ioutil.LogInfo{HttpUserAgent: userAgent3})
	handler.Output(limit)
	assert.Equal(t, 3, handler.countMap[userAgent1])
	assert.Equal(t, 2, handler.countMap[userAgent2])
	assert.Equal(t, 1, handler.countMap[userAgent3])
}

func TestMostFrequentResponseStatus(t *testing.T) {
	handler := NewMostFrequentStatusHandler()
	handler.Input(&ioutil.LogInfo{Status: responseStatus1, Request: uri1})
	handler.Input(&ioutil.LogInfo{Status: responseStatus2, Request: uri1})
	handler.Input(&ioutil.LogInfo{Status: responseStatus3, Request: uri1})
	handler.Input(&ioutil.LogInfo{Status: responseStatus1, Request: uri2})
	handler.Input(&ioutil.LogInfo{Status: responseStatus2, Request: uri2})
	handler.Input(&ioutil.LogInfo{Status: responseStatus3, Request: uri3})
	handler.Output(limit)
	assert.Equal(t, 2, handler.statusCountMap[responseStatus1])
	assert.Equal(t, 2, handler.statusCountMap[responseStatus2])
	assert.Equal(t, 2, handler.statusCountMap[responseStatus3])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus1][uri1])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus1][uri2])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus2][uri1])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus2][uri2])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus3][uri1])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus3][uri3])
}

func TestTopTimeMeanCostUris(t *testing.T) {
	handler := NewTopTimeMeanCostUrisHandler()
	handler.Input(&ioutil.LogInfo{Request: uri1, RequestTime: responseTime1})
	handler.Input(&ioutil.LogInfo{Request: uri1, RequestTime: responseTime2})
	handler.Input(&ioutil.LogInfo{Request: uri1, RequestTime: responseTime3})
	handler.Input(&ioutil.LogInfo{Request: uri2, RequestTime: responseTime1})
	handler.Input(&ioutil.LogInfo{Request: uri2, RequestTime: responseTime2})
	handler.Input(&ioutil.LogInfo{Request: uri3, RequestTime: responseTime1})
	handler.Output(limit)
	assert.Equal(t, []float64{responseTime1, responseTime2, responseTime3}, handler.timeCostListMap[uri1])
	assert.Equal(t, []float64{responseTime1, responseTime2}, handler.timeCostListMap[uri2])
	assert.Equal(t, []float64{responseTime1}, handler.timeCostListMap[uri3])
}

func TestTopTimePercentCostUris(t *testing.T) {
	handler := NewTopTimePercentCostUrisHandler(50)
	handler.Input(&ioutil.LogInfo{Request: uri1, RequestTime: responseTime1})
	handler.Input(&ioutil.LogInfo{Request: uri1, RequestTime: responseTime2})
	handler.Input(&ioutil.LogInfo{Request: uri1, RequestTime: responseTime3})
	handler.Input(&ioutil.LogInfo{Request: uri2, RequestTime: responseTime1})
	handler.Input(&ioutil.LogInfo{Request: uri2, RequestTime: responseTime2})
	handler.Input(&ioutil.LogInfo{Request: uri3, RequestTime: responseTime1})
	handler.Output(limit)
	assert.Equal(t, []float64{responseTime1, responseTime2, responseTime3}, handler.timeCostListMap[uri1])
	assert.Equal(t, []float64{responseTime1, responseTime2}, handler.timeCostListMap[uri2])
	assert.Equal(t, []float64{responseTime1}, handler.timeCostListMap[uri3])
}
