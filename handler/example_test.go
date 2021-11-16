package handler

import "github.com/fantasticmao/nginx-log-analyzer/parser"

func ExampleNewPvAndUvHandler() {
	handler := NewPvAndUvHandler()
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip2})
	handler.Input(&parser.LogInfo{RemoteAddr: ip2})
	handler.Input(&parser.LogInfo{RemoteAddr: ip3})
	handler.Output(limit)
	// Output:
	// PV: 6
	// UV: 3
}

func ExampleNewMostVisitedIpsHandler() {
	handler := NewMostVisitedFieldsHandler(AnalysisTypeVisitedIps)
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip1})
	handler.Input(&parser.LogInfo{RemoteAddr: ip2})
	handler.Input(&parser.LogInfo{RemoteAddr: ip2})
	handler.Input(&parser.LogInfo{RemoteAddr: ip3})
	handler.Output(limit)
	// Output:
	// "192.168.1.1" hits: 3
	// "192.168.1.2" hits: 2
	// "192.168.1.3" hits: 1
}

func ExampleNewMostVisitedUrisHandler() {
	handler := NewMostVisitedFieldsHandler(AnalysisTypeVisitedUris)
	handler.Input(&parser.LogInfo{Request: uri1})
	handler.Input(&parser.LogInfo{Request: uri1})
	handler.Input(&parser.LogInfo{Request: uri1})
	handler.Input(&parser.LogInfo{Request: uri2})
	handler.Input(&parser.LogInfo{Request: uri2})
	handler.Input(&parser.LogInfo{Request: uri3})
	handler.Output(limit)
	// Output:
	// "GET /name/Tom HTTP/2.0" hits: 3
	// "GET /name/Sam HTTP/2.0" hits: 2
	// "GET /name/Bob HTTP/2.0" hits: 1
}

func ExampleNewMostVisitedUserAgentsHandler() {
	handler := NewMostVisitedFieldsHandler(AnalysisTypeVisitedUserAgents)
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent1})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent2})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent2})
	handler.Input(&parser.LogInfo{HttpUserAgent: userAgent3})
	handler.Output(limit)
	// Output:
	// "iOS" hits: 3
	// "Android" hits: 2
	// "Windows" hits: 1
}

func ExampleNewMostVisitedLocationsHandler() {
	handler := NewMostVisitedLocationsHandler("../testdata/GeoLite2-City-Test.mmdb", limit)

	// see https://github.com/maxmind/MaxMind-DB/blob/main/source-data/GeoLite2-City-Test.json
	handler.Input(&parser.LogInfo{RemoteAddr: "175.16.199.0"}) // China -> Changchun
	handler.Input(&parser.LogInfo{RemoteAddr: "175.16.199.0"})
	handler.Input(&parser.LogInfo{RemoteAddr: "175.16.199.0"})
	handler.Input(&parser.LogInfo{RemoteAddr: "2.125.160.216"}) // United Kingdom -> Boxford
	handler.Input(&parser.LogInfo{RemoteAddr: "2.125.160.216"})
	handler.Input(&parser.LogInfo{RemoteAddr: "2001:218::"}) // Japan -> unknown
	handler.Output(limit)
	// Output:
	// [中国 China] hits: 3
	//   |--[长春 Changchun] hits: 3
	//   |  |--"175.16.199.0" hits: 3
	// [United Kingdom] hits: 2
	//   |--[Boxford] hits: 2
	//   |  |--"2.125.160.216" hits: 2
	// [日本 Japan] hits: 1
	//   |--[unknown] hits: 1
	//   |  |--"2001:218::" hits: 1
}

func ExampleNewMostFrequentStatusHandler() {
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
	// Output:
	// 200 hits: 3
	//   |--"GET /name/Tom HTTP/2.0" hits: 2
	//   |--"GET /name/Sam HTTP/2.0" hits: 1
	// 302 hits: 3
	//   |--"GET /name/Tom HTTP/2.0" hits: 2
	//   |--"GET /name/Sam HTTP/2.0" hits: 1
	// 404 hits: 3
	//   |--"GET /name/Tom HTTP/2.0" hits: 2
	//   |--"GET /name/Bob HTTP/2.0" hits: 1
}

func ExampleNewLargestAverageTimeUrisHandler() {
	handler := NewLargestAverageTimeUrisHandler()
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime1})
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime2})
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime3})
	handler.Input(&parser.LogInfo{Request: uri2, RequestTime: responseTime1})
	handler.Input(&parser.LogInfo{Request: uri2, RequestTime: responseTime2})
	handler.Input(&parser.LogInfo{Request: uri3, RequestTime: responseTime1})
	handler.Output(limit)
	// Output:
	// "GET /name/Tom HTTP/2.0" average response-time: 0.200
	// "GET /name/Sam HTTP/2.0" average response-time: 0.150
	// "GET /name/Bob HTTP/2.0" average response-time: 0.100
}

func ExampleNewLargestPercentTimeUrisHandler() {
	handler := NewLargestPercentTimeUrisHandler(30)
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime1})
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime2})
	handler.Input(&parser.LogInfo{Request: uri1, RequestTime: responseTime3})
	handler.Input(&parser.LogInfo{Request: uri2, RequestTime: responseTime2})
	handler.Input(&parser.LogInfo{Request: uri2, RequestTime: responseTime3})
	handler.Input(&parser.LogInfo{Request: uri3, RequestTime: responseTime3})
	handler.Output(limit)
	// Output:
	// "GET /name/Bob HTTP/2.0" P30.00 response-time: 0.300
	// "GET /name/Sam HTTP/2.0" P30.00 response-time: 0.200
	// "GET /name/Tom HTTP/2.0" P30.00 response-time: 0.100
}
