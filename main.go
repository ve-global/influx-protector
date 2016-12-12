package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/influxdata/influxdb/influxql"
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
	defaultPort           = ":8087"
	defaultPortUsage      = "default server port, ':8087', ':8086'..."
	defaultTarget         = "http://127.0.0.1:8086"
	defaultTargetUsage    = "default redirect url, 'http://127.0.0.1:8086'"
	defaultVerbose        = false
	defaultVerboseUsage   = "--verbose"
	defaultVersion        = false
	defaultVersionUsage   = "--version"
	defaultBuckets        = 1000
	defaultBucketsUsage   = "default buckets 1000"
	defaultSlowQuery      = 1000
	defaultSlowQueryUsage = "default slowquery time in milliseconds 1000"
)

func main() {
	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	target := flag.String("target", defaultTarget, defaultTargetUsage)
	verbose := flag.Bool("verbose", defaultVerbose, defaultVerboseUsage)
	vsn := flag.Bool("version", defaultVersion, defaultVersionUsage)
	maxbuckets := flag.Int("maxbuckets", defaultBuckets, defaultBucketsUsage)
	slowquery := flag.Int64("slowquery", defaultSlowQuery, defaultSlowQueryUsage)

	flag.Parse()

	if *vsn {
		fmt.Printf("influx-protector version %s", version.Version)
		fmt.Println()
		return
	}

	log.Printf("[INFO] Influx Protector Version %s", version.Version)
	log.Printf("[INFO] server will run on: %s", *port)
	log.Printf("[INFO] redirecting to: %s", *target)

	purl, _ := url.Parse(*target)
	proxy := httputil.NewSingleHostReverseProxy(purl)

	options := &rules.Options{
		MaxBuckets: *maxbuckets,
		SlowQuery:  *slowquery,
	}

	// server
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Influx-Protector-Version", version.Version)
		inputQuery := strings.TrimSpace(r.URL.Query().Get("q"))
		defer timeTrack(time.Now(), inputQuery, options)

		query, err := influxql.NewParser(strings.NewReader(inputQuery)).ParseStatement()

		if err != nil {
			writeError(inputQuery, err, w)
			return
		}

		if err := rules.RunRules(query, options); err != nil {
			writeError(inputQuery, err, w)
			return
		}

		if *verbose {
			log.Printf("[QUERY] %s", query.String())
		}

		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Influx-Protector-Version", version.Version)
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(*port, nil)
}

func writeError(rawQuery string, err error, w http.ResponseWriter) {
	log.Printf("[ERROR] %s ('%s')", err, rawQuery)

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

func timeTrack(start time.Time, query string, options *rules.Options) {
	elapsed := int64(time.Since(start).Seconds() * 1000)
	if elapsed > options.SlowQuery {
		log.Printf("[SLOWQUERY] '%s' took %dms", query, elapsed)
	}
}
