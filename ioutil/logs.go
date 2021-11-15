package ioutil

import (
	"fmt"
	"os"
	"time"
)

const timeLocalFormat = "02/Jan/2006:15:04:05 -0700"

func ParseTime(timeLocal string) time.Time {
	t, err := time.Parse(timeLocalFormat, timeLocal)
	if err != nil {
		Fatal("parse log time error: %v\n", err.Error())
	}
	return t
}

type LogInfo struct {
	RemoteAddr    string  `json:"remote_addr"`
	RemoteUser    string  `json:"remote_user"`
	TimeLocal     string  `json:"time_local"`
	Request       string  `json:"request"`
	Status        int     `json:"status"`
	BodyBytesSent int     `json:"body_bytes_sent"`
	HttpReferer   string  `json:"http_referer"`
	HttpUserAgent string  `json:"http_user_agent"`
	RequestTime   float64 `json:"request_time"`
}

func Fatal(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
