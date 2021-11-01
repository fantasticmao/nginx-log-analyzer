# Nginx-JSON-Log-Analyze

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

relate document: http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format

## Supported Statistical Indicators

- [x] PV and UV
- [x] Most visited URIs
- [x] Most visited IPs
- [x] Most visited User-Agents
- [x] Most time-cost URIs
- [ ] Most visited Cities
