package main

import (
	"flag"
	"fmt"
	"github.com/onidoru/FCFS_RR_Servers/pkg/queues"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	ArrivalRate = flag.Float64("l", 3.0, "Average arrival rate")
	ServiceRate = flag.Float64("u", 2, "Average service rate")
	N           = flag.Int("n", 1000, "Number of POST requests")
	HTTPAddr    = flag.String("http", "http://localhost:1337/work", "Address to to send HTTP requests")
)

func main() {
	flag.Parse()

	pD := *queues.PoissonDistribution(*ArrivalRate, *N)
	nD := *queues.NormalDistribution(*N, *ServiceRate, 0)
	http.PostForm(*HTTPAddr, url.Values{"limit": {strconv.Itoa(*N)}})
	for i, p := range pD {
		t := strconv.FormatFloat(nD[i], 'f', -1, 64) + "ms"
		fmt.Printf("POST: %v, \tNext delay: %v\n", t, p)
		http.PostForm(*HTTPAddr, url.Values{"delay": {t}})
		// fmt.Println(r)
		time.Sleep(time.Duration(p) * time.Millisecond)
	}
}
