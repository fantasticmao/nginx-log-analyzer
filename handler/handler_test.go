package handler

import (
	"github.com/fantasticmao/nginx-log-analyzer/parser"
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

func TestNewPvAndUvHandler(t *testing.T) {
	handler := NewPvAndUvHandler()
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip2})
	handler.Input(&parser.LogInfo{RemoteAddr: ip2})
	handler.Input(&parser.LogInfo{RemoteAddr: ip3})
	handler.Output(limit)

	assert.Equal(t, 6, handler.pv)
	assert.Equal(t, 3, handler.uv)
}

func TestNewMostVisitedFieldsHandler_ips(t *testing.T) {
	handler := NewMostVisitedFieldsHandler(AnalysisTypeVisitedIps)
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip2})
	handler.Input(&parser.LogInfo{RemoteAddr: ip2})
	handler.Input(&parser.LogInfo{RemoteAddr: ip3})
	handler.Output(limit)

	assert.Equal(t, 3, handler.countMap[ip1])
	assert.Equal(t, 2, handler.countMap[ip2])
	assert.Equal(t, 1, handler.countMap[ip3])
}

func TestNewMostVisitedFieldsHandler_uris(t *testing.T) {
	handler := NewMostVisitedFieldsHandler(AnalysisTypeVisitedUris)
	handler.Input(&parser.LogInfo{Request: uri1})
	handler.Input(&parser.LogInfo{Request: uri1})
	handler.Input(&parser.LogInfo{Request: uri1})
	handler.Input(&parser.LogInfo{Request: uri2})
	handler.Input(&parser.LogInfo{Request: uri2})
	handler.Input(&parser.LogInfo{Request: uri3})
	handler.Output(limit)

	assert.Equal(t, 3, handler.countMap[uri1])
	assert.Equal(t, 2, handler.countMap[uri2])
	assert.Equal(t, 1, handler.countMap[uri3])
}

func TestNewMostVisitedFieldsHandler_userAgents(t *testing.T) {
	handler := NewMostVisitedFieldsHandler(AnalysisTypeVisitedUserAgents)
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent2})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent2})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent3})
	handler.Output(limit)

	assert.Equal(t, 3, handler.countMap[userAgent1])
	assert.Equal(t, 2, handler.countMap[userAgent2])
	assert.Equal(t, 1, handler.countMap[userAgent3])
}

func TestNewMostVisitedLocationsHandler(t *testing.T) {
	handler := NewMostVisitedLocationsHandler("../testdata/GeoLite2-City-Test.mmdb", limit)
	assert.NotNil(t, handler.geoLite2Db)

	// see https://github.com/maxmind/MaxMind-DB/blob/main/source-data/GeoLite2-City-Test.json
	handler.Input(&parser.LogInfo{RemoteAddr: "175.16.199.0"}) // China -> Changchun
	handler.Input(&parser.LogInfo{RemoteAddr: "175.16.199.0"})
	handler.Input(&parser.LogInfo{RemoteAddr: "175.16.199.0"})
	handler.Input(&parser.LogInfo{RemoteAddr: "2.125.160.216"}) // United Kingdom -> Boxford
	handler.Input(&parser.LogInfo{RemoteAddr: "2.125.160.216"})
	handler.Input(&parser.LogInfo{RemoteAddr: "2001:218::"}) // Japan -> unknown
	handler.Output(limit)

	assert.Equal(t, 3, handler.countryCountMap["China"])
	assert.Equal(t, 2, handler.countryCountMap["United Kingdom"])
	assert.Equal(t, 1, handler.countryCountMap["Japan"])

	assert.Equal(t, 3, handler.countryCityCountMap["China"]["Changchun"])
	assert.Equal(t, 2, handler.countryCityCountMap["United Kingdom"]["Boxford"])
	assert.Equal(t, 1, handler.countryCityCountMap["Japan"]["unknown"])

	assert.Equal(t, 3, handler.countryCityIpCountMap["China"]["Changchun"]["175.16.199.0"])
	assert.Equal(t, 2, handler.countryCityIpCountMap["United Kingdom"]["Boxford"]["2.125.160.216"])
	assert.Equal(t, 1, handler.countryCityIpCountMap["Japan"]["unknown"]["2001:218::"])
}

func TestNewMostFrequentStatusHandler(t *testing.T) {
	handler := NewMostFrequentStatusHandler()
	handler.Input(&parser.LogInfo{Status: responseStatus1, Request: uri1})
	handler.Input(&parser.LogInfo{Status: responseStatus1, Request: uri1})
	handler.Input(&parser.LogInfo{Status: responseStatus2, Request: uri1})
	handler.Input(&parser.LogInfo{Status: responseStatus2, Request: uri1})
	handler.Input(&parser.LogInfo{Status: responseStatus3, Request: uri1})
	handler.Input(&parser.LogInfo{Status: responseStatus3, Request: uri1})
	handler.Input(&parser.LogInfo{Status: responseStatus1, Request: uri2})
	handler.Input(&parser.LogInfo{Status: responseStatus2, Request: uri2})
	handler.Input(&parser.LogInfo{Status: responseStatus3, Request: uri3})
	handler.Output(limit)

	assert.Equal(t, 3, handler.statusCountMap[responseStatus1])
	assert.Equal(t, 3, handler.statusCountMap[responseStatus2])
	assert.Equal(t, 3, handler.statusCountMap[responseStatus3])

	assert.Equal(t, 2, handler.statusUriCountMap[responseStatus1][uri1])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus1][uri2])
	assert.Equal(t, 2, handler.statusUriCountMap[responseStatus2][uri1])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus2][uri2])
	assert.Equal(t, 2, handler.statusUriCountMap[responseStatus3][uri1])
	assert.Equal(t, 1, handler.statusUriCountMap[responseStatus3][uri3])
}

func TestNewLargestAverageTimeUrisHandler(t *testing.T) {
	handler := NewLargestAverageTimeUrisHandler()
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime1})
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime2})
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime3})
	handler.Input(&parser.LogInfo{Request: uri2, RequestTime: responseTime1})
	handler.Input(&parser.LogInfo{Request: uri2, RequestTime: responseTime2})
	handler.Input(&parser.LogInfo{Request: uri3, RequestTime: responseTime1})
	handler.Output(limit)

	assert.Equal(t, []float64{responseTime1, responseTime2, responseTime3}, handler.timeCostListMap[uri1])
	assert.Equal(t, []float64{responseTime1, responseTime2}, handler.timeCostListMap[uri2])
	assert.Equal(t, []float64{responseTime1}, handler.timeCostListMap[uri3])
}

func TestNewLargestPercentTimeUrisHandler(t *testing.T) {
	handler := NewLargestPercentTimeUrisHandler(30)
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime1})
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime2})
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime3})
	handler.Input(&parser.LogInfo{Request: uri2, RequestTime: responseTime2})
	handler.Input(&parser.LogInfo{Request: uri2, RequestTime: responseTime3})
	handler.Input(&parser.LogInfo{Request: uri3, RequestTime: responseTime3})
	handler.Output(limit)

	assert.Equal(t, []float64{responseTime1, responseTime2, responseTime3}, handler.timeCostListMap[uri1])
	assert.Equal(t, []float64{responseTime2, responseTime3}, handler.timeCostListMap[uri2])
	assert.Equal(t, []float64{responseTime3}, handler.timeCostListMap[uri3])
}
