# Status Board

### build
`make`

### help
`./status-board --help`

### run
`/status-board --port=8080 --sites_path=/path/to/sites.txt --metrics --timeout=5 --check_rate=60`

## Check status
```
GET /status/min
GET /status/max
GET /status/random
GET /status/site/{site_name}
```

## Metrics
```
GET /metrics
GET /metrics/{site_name}
```
