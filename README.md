Influxdb Protector
---

[ ![Build Status](https://travis-ci.org/ve-interactive/influx-protector.svg?branch=master)](https://travis-ci.org/ve-interactive/influx-protector)

Inspired by [Protector](https://github.com/trivago/Protector). A proxy written in Golang that will protect influxdb from silly or dangerous queries.

- Queries loading lots of data points (threshold is configurable)
- Queries dropping, altering or otherwise messing with the database
- Show Measurements queries
- Queries with short source regexes (i.e. likely to match lots of series)
- Queries without a where clause


## Options

- target: target server `--target http://127.0.0.1:8086`
- port: port to run on `--port 8087`
- verbose: log all queries `--verbose`
- maxbuckets: max number of data points for a single query: `--maxbuckets 1400`
- slowqueries: set slow queries threshold (milliseconds): `--slowqueries 1000`

## Supported endpoints

- `/query?q=...`
- `/ping`
