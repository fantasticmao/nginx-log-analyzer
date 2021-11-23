package handler

import (
	"math/rand"
	"testing"
)

var (
	// see https://github.com/maxmind/MaxMind-DB/blob/main/source-data/GeoLite2-City-Test.json
	ips = []string{
		"2.125.160.216", "67.43.156.0", "81.2.69.142",
		"81.2.69.144", "89.160.20.112", "175.16.199.0",
		"2001:218::", "2001:252::", "2001:230::",
	}
)

func BenchmarkQueryIpLocation(b *testing.B) {
	handler := NewMostVisitedLocationsHandler("../testdata/GeoLite2-City-Test.mmdb", limit)

	for i, l := 0, len(ips); i < b.N; i++ {
		ip := ips[rand.Intn(l)]
		_, _ = handler.queryIpLocation(ip)
	}
}

func BenchmarkCachedQueryIpLocation(b *testing.B) {
	handler := NewMostVisitedLocationsHandler("../testdata/GeoLite2-City-Test.mmdb", limit)

	for i, l := 0, len(ips); i < b.N; i++ {
		ip := ips[rand.Intn(l)]
		_, _ = handler.cachedQueryIpLocation(ip)
	}
}
