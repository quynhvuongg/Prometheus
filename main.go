package main

import (
	"flag"
	"log"
	"net/http"
    "fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"crypto/tls"
	"github.com/oliveagle/jsonpath"
	"io/ioutil"
	"encoding/json"
)

var addr = flag.String("listen-address", ":9116", "The address to listen on for HTTP requests.")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>Json Exporter</title></head>
            <body>
            <h1>Json Exporter</h1>
            <p><a href="/probe">Run a probe</a></p>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})
	flag.Parse()
	http.HandleFunc("/probe", probeHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func probeHandler(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	target := params.Get("http://localhost:9300/api/v1/silences")
	if target == "" {
		http.Error(w, "Target parameter is missing", 400)
		return
	}

	
	// lookuppath := params.Get("jsonpath")
	// if target == "" {
	// 	http.Error(w, "The JsonPath to lookup", 400)
	// 	return
	// }
    // var i int 0
	// lookuppath := params.Get("$data.id[%v]",i)
	// if target == "" {
	// 	http.Error(w, "The JsonPath to lookup", 400)
	// 	return
	// }


    }
	// probeSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "probe_success",
	// 	Help: "Displays whether or not the probe was a success",
	// })
	// probeDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "probe_duration_seconds",
	// 	Help: "Returns how long the probe took to complete in seconds",
	// })
	// valueGauge := prometheus.NewGauge(
	// 	prometheus.GaugeOpts{
	// 		Name:	"value",
	// 		Help:	"Retrieved value",
	// 	},
	// )

	ID := prometheus.NewUntyped(
		prometheus.UntypedOpts{
			Name:	"ID",
			Help:	"ID silences",
		},
	)

	Creator := prometheus.NewUntyped(
		prometheus.UntypedOpts{
			Name:	"Creator",
			Help:	"CreatedBy",
		},
	)

	StartsAt := prometheus.NewUntyped(
		prometheus.NewUntyped{
			Name:	"StartsAt",
			Help:	"Starts to create silences",
		},
	)

	EndsAt := prometheus.NewUntyped(
		prometheus.UntypedOpts{
			Name:	"EndsAt",
			Help:	"Ends to create silences",
		},
	)

	State := prometheus.NewUntyped(
		prometheus.UntypedOpts{
			Name:	"State",
			Help:	"State of silences",
		},
	)


	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)
	registry.MustRegister(probeDurationGauge)
	registry.MustRegister(ID)
	registry.MustRegister(Creator)
	registry.MustRegister(StartsAt)
	registry.MustRegister(EndsAt)
	registry.MustRegister(State)
	


	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(target)
	if err != nil {
		log.Fatal(err)

	} else {
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		
		var json_data interface{}
		json.Unmarshal([]byte(bytes), &json_data)
		for k:=0;k< 50; k++ {
		a, err := jsonpath.JsonPathLookup(json_data,"$.data.id[%v]",k)
		if err != nil {
			http.Error(w, "Jsonpath not found", http.StatusNotFound)
			return
			}
		log.Printf("Found value %v", a)

		b, err := jsonpath.JsonPathLookup(json_data,"$.data.createdBy[%v]",k)
		if err != nil {
			http.Error(w, "Jsonpath not found", http.StatusNotFound)
			return
			}
		log.Printf("Found value %v", b)

		c, err := jsonpath.JsonPathLookup(json_data,"$.data.startsAt[%v]",k)
		if err != nil {
			http.Error(w, "Jsonpath not found", http.StatusNotFound)
			return
			}
		log.Printf("Found value %v", c)

		d, err := jsonpath.JsonPathLookup(json_data,"$.data.endsAt[%v]",k)
		if err != nil {
			http.Error(w, "Jsonpath not found", http.StatusNotFound)
			return
			}
		log.Printf("Found value %v", d)
		e, err := jsonpath.JsonPathLookup(json_data,"$.data.status[%v].state",k)
		if err != nil {
			http.Error(w, "Jsonpath not found", http.StatusNotFound)
			return
			}
		log.Printf("Found value %v", e)

        ID.Set(a)
		Creator.Set(b)
		StartsAt.Set(c)
		EndsAt.Set(d)
		State.Set(e)
		}
	}	

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}