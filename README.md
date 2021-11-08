# Nginx-JSON-Log-Analyzer

[![Actions Status](https://github.com/fantasticmao/nginx-json-log-analyzer/workflows/ci/badge.svg)](https://github.com/fantasticmao/nginx-json-log-analyzer/actions)
[![codecov](https://codecov.io/gh/fantasticmao/nginx-json-log-analyzer/branch/main/graph/badge.svg)](https://codecov.io/gh/fantasticmao/nginx-json-log-analyzer)
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
access_log /path/to/access.json.log json_log;
```

The `log_format` directive can only appear in the `http` context.

The `access_log` directive should use the log format declared above, and you can make multiple `access_log`s at the same time without deleting the original configuration. e.g.

```text
access_log /path/to/access.log;
access_log /path/to/access.json.log json_log;
```

Related document: http://nginx.org/en/docs/http/ngx_http_log_module.html

## Supported Statistical Indicators

| Supported | Analyze Type | Statistical Indicators                                                       | Required Fields or Libraries                                                                                                                                     |
| --------- | ------------ | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ✅        | 0            | PV and UV                                                                    | $remote_addr                                                                                                                                                     |
| ✅        | 1            | Most visited IPs                                                             | $remote_addr                                                                                                                                                     |
| ✅        | 2            | Most visited URIs                                                            | $request                                                                                                                                                         |
| ✅        | 3            | Most visited User-Agents                                                     | $http_user_agent                                                                                                                                                 |
| ✅        | 4            | Most visited user countries and cities                                       | $remote_addr, MaxMind [GeoIP2](https://www.maxmind.com/en/geoip2-city) or [GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) City Database |
| ✅        | 5            | Most frequent response status                                                | $status, $request                                                                                                                                                |
| ✅        | 6            | Top mean response-time URIs                                                  | $request, $request_time                                                                                                                                          |
| ✅        | 7            | Top percentile response-time URIs, e.g. p1(min), p50(median), p95, p100(max) | $request, $request_time                                                                                                                                          |
