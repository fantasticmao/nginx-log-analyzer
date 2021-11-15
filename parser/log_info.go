package parser

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
