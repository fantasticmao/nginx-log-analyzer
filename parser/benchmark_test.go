package parser

import "testing"

func BenchmarkJsonParser(b *testing.B) {
	p := NewJsonParser()
	for i := 0; i < b.N; i++ {
		p.ParseLog(jsonLog)
	}
}

func BenchmarkCombinedParser(b *testing.B) {
	p := NewCombinedParser()
	for i := 0; i < b.N; i++ {
		p.ParseLog(combinedLog)
	}
}
