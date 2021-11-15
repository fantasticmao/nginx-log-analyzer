package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTime(t *testing.T) {
	datetime := ParseTime("01/Nov/2021:00:00:00 +0800")
	assert.NotNil(t, datetime)
	assert.Equal(t, int64(1635696000000), datetime.UnixMilli())
}

func TestParseLogJson(t *testing.T) {
	log := "{\"remote_addr\":\"66.102.6.200\",\"time_local\":\"15/Nov/2021:13:44:10 +0800\",\"request\":\"GET / HTTP/1.1\",\"status\":200,\"body_bytes_sent\":1603,\"http_user_agent\":\"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.75 Safari/537.36 Google Favicon\",\"request_time\":0.20}\n"
	logInfo := NewJsonParser().ParseLog([]byte(log))
	assert.NotNil(t, logInfo)
	assert.Equal(t, "66.102.6.200", logInfo.RemoteAddr)
	assert.Equal(t, "", logInfo.RemoteUser)
	assert.Equal(t, "15/Nov/2021:13:44:10 +0800", logInfo.TimeLocal)
	assert.Equal(t, "GET / HTTP/1.1", logInfo.Request)
	assert.Equal(t, 200, logInfo.Status)
	assert.Equal(t, 1603, logInfo.BodyBytesSent)
	assert.Equal(t, "", logInfo.HttpReferer)
	assert.Equal(t, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.75 Safari/537.36 Google Favicon", logInfo.HttpUserAgent)
	assert.Equal(t, 0.20, logInfo.RequestTime)
}

func TestParseLogCombined(t *testing.T) {
	log := "66.102.6.200 - - [15/Nov/2021:13:44:10 +0800] \"GET / HTTP/1.1\" 200 1603 \"-\" \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.75 Safari/537.36 Google Favicon\"\n"
	logInfo := NewCombinedParser().ParseLog([]byte(log))
	assert.NotNil(t, logInfo)
	assert.Equal(t, "66.102.6.200", logInfo.RemoteAddr)
	assert.Equal(t, "-", logInfo.RemoteUser)
	assert.Equal(t, "15/Nov/2021:13:44:10 +0800", logInfo.TimeLocal)
	assert.Equal(t, "GET / HTTP/1.1", logInfo.Request)
	assert.Equal(t, 200, logInfo.Status)
	assert.Equal(t, 1603, logInfo.BodyBytesSent)
	assert.Equal(t, "-", logInfo.HttpReferer)
	assert.Equal(t, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.75 Safari/537.36 Google Favicon", logInfo.HttpUserAgent)
}
