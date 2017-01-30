package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/influxdata/influxdb/influxql"
	"github.com/ve-interactive/influx-protector/logger"
	"github.com/ve-interactive/influx-protector/rules"
	"github.com/ve-interactive/influx-protector/version"
)

type errorResponses struct {
	Results []*errorResponse `json:"results"`
}

type errorResponse struct {
	Error string `json:"error"`
}

const (
	defaultPort               = ":8087"
	defaultPortUsage          = "default server port, ':8087', ':8086'..."
	defaultTarget             = "http://127.0.0.1:8086"
	defaultTargetUsage        = "default redirect url, 'http://127.0.0.1:8086'"
	defaultVerbose            = false
	defaultVerboseUsage       = "--verbose"
	defaultVersion            = false
	defaultVersionUsage       = "--version"
	defaultBuckets            = 1000
	defaultBucketsUsage       = "default buckets 1000"
	defaultSlowQuery          = 1000
	defaultSlowQueryUsage     = "default slowquery time in milliseconds 1000"
	defaultStatsdAddress      = "localhost:8125"
	defaultStatsdAddressUsage = "<statsd_address>:<port>"
)

func main() {
	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	target := flag.String("target", defaultTarget, defaultTargetUsage)
	verbose := flag.Bool("verbose", defaultVerbose, defaultVerboseUsage)
	vsn := flag.Bool("version", defaultVersion, defaultVersionUsage)
	maxbuckets := flag.Int("maxbuckets", defaultBuckets, defaultBucketsUsage)
	slowquery := flag.Int64("slowquery", defaultSlowQuery, defaultSlowQueryUsage)
	statsdAddress := flag.String("statsdaddress", defaultStatsdAddress, defaultStatsdAddressUsage)

	flag.Parse()

	if *vsn {
		fmt.Printf("influx-protector version %s", version.Version)
		fmt.Println()
		return
	}

	logger := logger.NewLogger(*verbose, *slowquery)
	sink, _ := metrics.NewStatsdSink(*statsdAddress)
	metrics.NewGlobal(metrics.DefaultConfig("influx-protector"), sink)

	logger.Info("Influx Protector Version %s", version.Version)
	logger.Info("server will run on: %s", *port)
	logger.Info("redirecting to: %s", *target)

	purl, _ := url.Parse(*target)
	proxy := httputil.NewSingleHostReverseProxy(purl)

	options := &rules.Options{
		MaxBuckets: *maxbuckets,
		SlowQuery:  *slowquery,
	}

	// server
	http.HandleFunc("/query", queryFunc(logger, options, proxy))

	http.HandleFunc("/ping", pingFunc(logger, options, proxy))

	http.HandleFunc("/write", writeFunc(logger, options, proxy))

	http.ListenAndServe(*port, nil)
}

func generateErrorResp(err error, w http.ResponseWriter) {
	body, jsErr := json.Marshal(&errorResponses{
		Results: []*errorResponse{&errorResponse{
			Error: err.Error(),
		}},
	})

	if jsErr != nil {
		http.Error(w, jsErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Influx-Protector-Version", version.Version)
}

func queryFunc(logger *logger.Logger, options *rules.Options, proxy *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w)
		inputQuery := strings.TrimSpace(r.URL.Query().Get("q"))
		defer logger.Query(time.Now(), inputQuery, options)
		defer metrics.MeasureSince([]string{"queries", "timing"}, time.Now())

		query, err := influxql.NewParser(strings.NewReader(inputQuery)).ParseStatement()

		if err != nil {
			logger.Error(inputQuery, err)
			metrics.IncrCounter([]string{"queries", "malformed"}, 1)
			generateErrorResp(err, w)
			return
		}

		if err := rules.RunRules(query, options); err != nil {
			logger.Error(inputQuery, err)
			metrics.IncrCounter([]string{"queries", "blocked"}, 1)
			generateErrorResp(err, w)
			return
		}

		metrics.IncrCounter([]string{"queries", "accepted"}, 1)
		proxy.ServeHTTP(w, r)
	}
}

func pingFunc(logger *logger.Logger, options *rules.Options, proxy *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w)
		proxy.ServeHTTP(w, r)
	}
}

func writeFunc(logger *logger.Logger, options *rules.Options, proxy *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w)
		defer metrics.MeasureSince([]string{"writes", "timing"}, time.Now())

		buf, _ := ioutil.ReadAll(r.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))

		body := string(buf[:len(buf)])
		lines := strings.Split(body, "\n")
		r.Body = rdr1

		datapoints := float32(len(lines))
		metrics.IncrCounter([]string{"writes", "points"}, datapoints)
		metrics.SetGauge([]string{"writes", "batchsize"}, datapoints)

		proxy.ServeHTTP(w, r)
	}
}
