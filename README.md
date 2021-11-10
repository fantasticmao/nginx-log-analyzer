# Nginx-JSON-Log-Analyzer

[![Actions Status](https://github.com/fantasticmao/nginx-json-log-analyzer/workflows/ci/badge.svg)](https://github.com/fantasticmao/nginx-json-log-analyzer/actions)
[![codecov](https://codecov.io/gh/fantasticmao/nginx-json-log-analyzer/branch/main/graph/badge.svg)](https://codecov.io/gh/fantasticmao/nginx-json-log-analyzer)
![Go Version](https://img.shields.io/github/go-mod/go-version/fantasticmao/nginx-json-log-analyzer)
[![Go Report Card](https://goreportcard.com/badge/github.com/fantasticmao/nginx-json-log-analyzer)](https://goreportcard.com/report/github.com/fantasticmao/nginx-json-log-analyzer)
[![Release](https://img.shields.io/github/v/release/fantasticmao/nginx-json-log-analyzer)](https://github.com/fantasticmao/nginx-json-log-analyzer/releases)
[![License](https://img.shields.io/github/license/fantasticmao/nginx-json-log-analyzer)](https://github.com/fantasticmao/nginx-json-log-analyzer/blob/main/LICENSE)

README [English](README.md) | [中文](README_ZH.md)

## What is it

### Advantages compared to [GoAccess](https://goaccess.io/)

### Advantages compared to [ELK](https://www.elastic.co/cn/what-is/elk-stack)

## Quick start

### Installation

#### GeoIP2 and GeoLite2

### Configure Nginx

Nginx-JSON-Log-Analyzer can only parse Nginx access logs in JSON format, so you need to add the following `log_format` and `access_log` directives in the Nginx configuration:

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

- The `log_format` directive can only appear in the `http` context;
- The `access_log` directive could appear in the `http`, `server`, `location` context, and should use the `log_format` declared above;
- You can make multiple `access_log`s at the same time without deleting the original configuration. e.g.
    ```text
    access_log /path/to/access.log;
    access_log /path/to/access.json.log json_log;
    ```

Related document: http://nginx.org/en/docs/http/ngx_http_log_module.html

### Command line options

#### show version -v

The `-v` options show Nginx-JSON-Log-Analyzer's build version, build time, and Git Commit at build time.

#### specify the configuration directory -d

The `-d` option specify the configuration directory that Nginx-JSON-Log-Analyzer required at runtime, the default value is `${HOME}/.config/nginx-json-log-analyzer/`.

#### specify the analysis type -t

The `-t` option specify the type of this analysis, the analysis type and corresponding statistical indicators are as follows:

| Supported | Analysis Type `-t` | Statistical Indicators                                                       | Required Fields or Libraries                                                                                                                                     |
| --------- | ------------------ | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ✅        | 0                  | PV and UV                                                                    | $remote_addr                                                                                                                                                     |
| ✅        | 1                  | Most visited IPs                                                             | $remote_addr                                                                                                                                                     |
| ✅        | 2                  | Most visited URIs                                                            | $request                                                                                                                                                         |
| ✅        | 3                  | Most visited User-Agents                                                     | $http_user_agent                                                                                                                                                 |
| ✅        | 4                  | Most visited user countries and cities                                       | $remote_addr, MaxMind [GeoIP2](https://www.maxmind.com/en/geoip2-city) or [GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) City Database |
| ✅        | 5                  | Most frequent response status                                                | $status, $request                                                                                                                                                |
| ✅        | 6                  | Top mean response-time URIs                                                  | $request, $request_time                                                                                                                                          |
| ✅        | 7                  | Top percentile response-time URIs, e.g. p1(min), p50(median), p95, p100(max) | $request, $request_time                                                                                                                                          |

#### limit the analysis start and end time -ta -tb

`-ta` and `-tb` options are used to filter logs based on the request time, ta is the abbreviation of time after, tb is the abbreviation of time before.

`-ta` and `-tb` options required the $time_iso8601 field in `log_format` directive of Nginx configuration.

#### limit the output lines number -n -n2

`-n` and `-n2` options are used to limit the number of output lines of Nginx-JSON-Log-Analyzer, `-n2` option only works in `-t 4` mode.

#### specify the percentile value -p

The `-p` option specify the percentile value in the `-t 7` mode, the default value is 95.

### Usages

## FQA

## License

GeoLite2 Database [License](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data#license)

Nginx-JSON-Log-Analyzer [License](https://github.com/fantasticmao/nginx-json-log-analyzer/blob/main/LICENSE)

Copyright (c) 2021 fantasticmao
