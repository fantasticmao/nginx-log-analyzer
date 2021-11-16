package handler

import (
	"math/rand"
	"net"
	"testing"
)

func BenchmarkQueryIpLocation(b *testing.B) {
	handler := NewMostVisitedLocationsHandler("../testdata/GeoLite2-City-Test.mmdb", limit)

	// see https://github.com/maxmind/MaxMind-DB/blob/main/source-data/GeoLite2-City-Test.json
	ips := []net.IP{
		net.ParseIP("2.125.160.216"), net.ParseIP("67.43.156.0"), net.ParseIP("81.2.69.142"),
		net.ParseIP("81.2.69.144"), net.ParseIP("89.160.20.112"), net.ParseIP("175.16.199.0"),
		net.ParseIP("2001:218::"), net.ParseIP("2001:252::"), net.ParseIP("2001:230::"),
	}
	for i, l := 0, len(ips); i < b.N; i++ {
		ip := ips[rand.Intn(l)]
		_, _ = handler.queryIpLocation(ip)
	}
}
