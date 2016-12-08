package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/influxdata/influxdb/influxql"
	"github.com/ve-interactive/influx-protector/rules"
	"github.com/ve-interactive/influx-protector/version"
)

type ErrorResponses struct {
	Results []*ErrorResponse `json:"results"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	// come constants and usage helper
	const (
		defaultPort         = ":8087"
		defaultPortUsage    = "default server port, ':8087', ':8086'..."
		defaultTarget       = "http://127.0.0.1:8086"
		defaultTargetUsage  = "default redirect url, 'http://127.0.0.1:8086'"
		defaultVerbose      = false
		defaultVerboseUsage = "--verbose"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	target := flag.String("target", defaultTarget, defaultTargetUsage)
	verbose := flag.Bool("verbose", defaultVerbose, defaultVerboseUsage)

	flag.Parse()

	log.Printf("[INFO] Influx Protector Version %s", version.Version)
	log.Printf("[INFO] server will run on: %s", *port)
	log.Printf("[INFO] redirecting to: %s", *target)

	purl, _ := url.Parse(*target)
	proxy := httputil.NewSingleHostReverseProxy(purl)

	// server
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Influx-Protector-Version", version.Version)
		inputQuery := strings.TrimSpace(r.URL.Query().Get("q"))
		query, err := influxql.NewParser(strings.NewReader(inputQuery)).ParseStatement()

		if err != nil {
			writeError(inputQuery, err, w)
			return
		}

		if err := rules.RunRules(query); err != nil {
			writeError(inputQuery, err, w)
			return
		}

		if *verbose {
			log.Printf("[QUERY] %s", query.String())
		}

		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(*port, nil)
}

func writeError(rawQuery string, err error, w http.ResponseWriter) {
	log.Printf("[ERROR] %s ('%s')", err, rawQuery)

	body, jsErr := json.Marshal(&ErrorResponses{
		Results: []*ErrorResponse{&ErrorResponse{
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
