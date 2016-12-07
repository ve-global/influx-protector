Influxdb Protector
---

Inspired by [Protector](https://github.com/trivago/Protector). A proxy written in Golang that will protect influxdb from silly or dangerous queries.

- Queries loading lots of data points (threshold is configurable)
- Queries dropping, altering or otherwise messing with the database
- Show Measurements queries
- Queries with short source regexes (i.e. likely to match lots of series)
- Queries loading raw data (i.e. without aggregation)
- Queries without a where clause
