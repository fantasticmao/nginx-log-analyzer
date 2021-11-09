# Nginx-JSON-Log-Analyzer

[![Actions Status](https://github.com/fantasticmao/nginx-json-log-analyzer/workflows/ci/badge.svg)](https://github.com/fantasticmao/nginx-json-log-analyzer/actions)
[![codecov](https://codecov.io/gh/fantasticmao/nginx-json-log-analyzer/branch/main/graph/badge.svg)](https://codecov.io/gh/fantasticmao/nginx-json-log-analyzer)
![Go Version](https://img.shields.io/github/go-mod/go-version/fantasticmao/nginx-json-log-analyzer)
[![Go Report Card](https://goreportcard.com/badge/github.com/fantasticmao/nginx-json-log-analyzer)](https://goreportcard.com/report/github.com/fantasticmao/nginx-json-log-analyzer)
[![License](https://img.shields.io/github/license/fantasticmao/nginx-json-log-analyzer)](https://github.com/fantasticmao/nginx-json-log-analyzer/blob/main/LICENSE)

[README](README.md) | [看我看我](README_ZH.md)

## 这是什么

Nginx-JSON-Log-Analyzer 是一个轻量的（简陋的）对 JSON 格式日志的分析工具，用于满足我自己对 Nginx 访问日志的分析需求。

Nginx-JSON-Log-Analyzer 基于 Go 语言来编写，运行时只需一个 2MB 左右的可执行文件，目前支持的功能特性如下：

- [x] 基于请求时间过滤数据
- [x] 支持同时解析多个文件
- [x] 支持解析 .gz 压缩文件
- [x] 支持多种 [统计指标](#统计指标)

### 和 [GoAccess](https://goaccess.io/) 相比有什么优势

在开发 Nginx-JSON-Log-Analyzer 之前，我不知道 GoAccess 的存在，不然可能就不会有这个仓库了。

GoAccess 是一个优秀的实时 web 日志解析工具，比 Nginx-JSON-Log-Analyzer 更好用、更强大。不过据我所知，GoAccess 似乎不支持读取 .gz 格式的压缩文件，也不支持按百分位统计 URI 的响应时间。

### 和 [ELK](https://www.elastic.co/cn/what-is/elk-stack) 相比有什么优势

ELK 虽然功能强大，但安装和配置比较麻烦，对机器性能也有一定要求。Nginx-JSON-Log-Analyzer 更加轻量，使用起来更加简单，适用于一些简单的日志分析场景。

## 快速开始

### 下载安装

在 Nginx-JSON-Log-Analyzer 的 GitHub [Release](https://github.com/fantasticmao/nginx-json-log-analyzer/releases) 页面中，下载对应平台的二进制可执行文件即可。

#### GeoIP2 和 GeoLite2

[GeoIP2](https://www.maxmind.com/en/geoip2-city) 是商业版的 IP 地理定位的数据库，需要付费才能使用。[GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) 是免费版和低精度版的 GeoIP2，以 [署名-相同方式共享 4.0 国际](https://creativecommons.org/licenses/by-sa/4.0/deed.zh) 许可证发行，在 [MaxMind](https://www.maxmind.com/en/accounts/current/geoip/downloads) 官网登录即可下载。

在使用 Nginx-JSON-Log-Analyzer 时，如果需要解析 IP 的地理位置（即使用 `-t 4` 模式），则需要额外下载 GeoIP2 或 GeoLite2 的城市数据库文件，保存至默认配置目录 `${HOME}/.config/nginx-json-log-analyzer/` 中的 `City.mmdb` 文件：

```shell
~$ mkdir ${HOME}/.config/nginx-json-log-analyzer
~$ tar -xzvf GeoLite2-City_20211102.tar.gz
~$ cp GeoLite2-City_20211102/GeoLite2-City.mmdb ${HOME}/.config/nginx-json-log-analyzer/City.mmdb
```

### Nginx 配置

Nginx-JSON-Log-Analyzer 只能解析 JSON 格式的 Nginx 访问日志，因此需要在 Nginx 配置文件中添加如下的配置：

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

- `log_format` 指令只能出现在 `http` 上下文中。
- `access_log` 指令可以出现在 `http`、`server`、`location` 等上下文中，需要指定使用如上声明的 `log_format`， 并且你可以同时使用多个 `access_log`，而不用删除原先已有的配置。例如：
    ```text
    access_log /path/to/access.log;
    access_log /path/to/access.json.log json_log;
    ```

相关文档: http://nginx.org/en/docs/http/ngx_http_log_module.html

### 命令行选项

TODO

### 统计指标

| 是否支持 | 分析类型 | 统计指标                                                           | 需要的字段或者依赖                                                                                                                                               |
| -------- | -------- | ------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ✅       | 0        | PV 和 UV                                                           | $remote_addr                                                                                                                                                     |
| ✅       | 1        | 访问最多的 IP                                                      | $remote_addr                                                                                                                                                     |
| ✅       | 2        | 访问最多的 URI                                                     | $request                                                                                                                                                         |
| ✅       | 3        | 访问最多的 User-Agent                                              | $http_user_agent                                                                                                                                                 |
| ✅       | 4        | 访问最多的国家和城市                                               | $remote_addr, MaxMind [GeoIP2](https://www.maxmind.com/en/geoip2-city) or [GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) City Database |
| ✅       | 5        | 频率最高的响应状态码                                               | $status, $request                                                                                                                                                |
| ✅       | 6        | 最大 URI 平均响应时间                                              | $request, $request_time                                                                                                                                          |
| ✅       | 7        | 最大 URI 百分位响应时间，例如 p1(最小), p50(中位), p95, p100(最大) | $request, $request_time                                                                                                                                          |

### 使用示例

#### 基于请求时间过滤数据

#### 同时解析多个文件

#### 解析 .gz 压缩文件

####

## FAQ

Q: 未来是否会支持实时解析？

A: 不会支持。如果想要实时解析，建议使用 GoAccess、ELK、Grafana + 时序数据库之类的方案。

## 版权声明

GeoLite2 Database [版权声明](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data#license)

Nginx-JSON-Log-Analyzer [版权声明](https://github.com/fantasticmao/nginx-json-log-analyzer/blob/main/LICENSE)

Copyright (c) 2021 fantasticmao