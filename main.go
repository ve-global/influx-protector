package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/influxdata/influxdb/influxql"
	"github.com/ve-interactive/influx-protector/rules"
)

func main() {
	// come constants and usage helper
	const (
		defaultPort        = ":8087"
		defaultPortUsage   = "default server port, ':8087', ':8086'..."
		defaultTarget      = "http://127.0.0.1:8086"
		defaultTargetUsage = "default redirect url, 'http://127.0.0.1:8086'"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	target := flag.String("target", defaultTarget, defaultTargetUsage)

	flag.Parse()

	log.Printf("server will run on : %s", *port)
	log.Printf("redirecting to :%s", *target)

	purl, _ := url.Parse(*target)
	proxy := httputil.NewSingleHostReverseProxy(purl)

	// server
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-GoProxy", "GoProxy")

		inputQuery := strings.TrimSpace(r.URL.Query().Get("q"))
		query, err := influxql.NewParser(strings.NewReader(inputQuery)).ParseStatement()

		if err != nil {
			returnError(err, w)
			return
		}

		log.Println(query.String())

		if err := rules.RunRules(query); err != nil {
			returnError(err, w)
			return
		}

		// call to magic method from ReverseProxy object
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(*port, nil)
}

func returnError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(400)
}
