package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonLogParser(t *testing.T) {
	log := "{\"time_local\":\"01/Nov/2021:00:00:00 +0800\",\"remote_addr\":\"192.168.1.1\",\"request_time\":0.010,\"request\":\"GET /name/Tom HTTP/2.0\",\"status\":200,\"body_bytes_sent\":100,\"http_user_agent\":\"iOS\"}"
	logInfo := JsonLogParser([]byte(log))
	assert.NotNil(t, logInfo)
	assert.Equal(t, "01/Nov/2021:00:00:00 +0800", logInfo.TimeLocal)
	assert.Equal(t, "192.168.1.1", logInfo.RemoteAddr)
	assert.Equal(t, 0.01, logInfo.RequestTime)
	assert.Equal(t, "GET /name/Tom HTTP/2.0", logInfo.Request)
	assert.Equal(t, 200, logInfo.Status)
	assert.Equal(t, 100, logInfo.BodyBytesSent)
	assert.Equal(t, "iOS", logInfo.HttpUserAgent)
}
