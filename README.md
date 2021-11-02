# Nginx-JSON-Log-Analyzer

![Go Version](https://img.shields.io/github/go-mod/go-version/fantasticmao/nginx-json-log-analyzer)
[![Go Report Card](https://goreportcard.com/badge/github.com/fantasticmao/nginx-json-log-analyzer)](https://goreportcard.com/report/github.com/fantasticmao/nginx-json-log-analyzer)
[![License](https://img.shields.io/github/license/fantasticmao/nginx-json-log-analyzer)](https://github.com/fantasticmao/nginx-json-log-analyzer/blob/main/LICENSE)

## Nginx Configuration

```text
log_format json_log escape=json '{"time_iso8601":"$time_iso8601",'
                                '"remote_addr":"$remote_addr",'
                                '"request_time":$request_time,'
                                '"request":"$request",'
                                '"status":$status,'
                                '"body_bytes_sent":$body_bytes_sent,'
                                '"http_user_agent":"$http_user_agent"}';
access_log /path/to/access.log json_log
```

Related document: http://nginx.org/en/docs/http/ngx_http_log_module.html

## Supported Statistical Indicators

- [x] PV and UV
- [x] Most visited IPs
- [x] Most visited URIs
- [x] Most visited User-Agents
- [ ] Most visited Countries
- [ ] Most visited Cities
- [x] Top Average time-cost URIs
- [ ] Top 99% time-cost URIs
- [ ] Top 95% time-cost URIs
