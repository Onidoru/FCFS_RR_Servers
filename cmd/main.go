package main

import (
	"flag"
	"fmt"
	model "github.com/onidoru/FCFS_RR_Servers/pkg/queues"
	"net/http"
)

var (
	NFCFS    = flag.Int("fcfs", 1, "Number of FCFS servers")
	NRR      = flag.Int("rr", 1, "Number of RR servers")
	HRR      = flag.Float64("hrr", 1.1, "H for RR servers")
	T1       = flag.Float64("t1", 0.75, "T1 relevance measure")
	T2       = flag.Float64("t2", 1.5, "T2 relevance measure")
	HTTPAddr = flag.String("http", "127.0.0.1:1337", "Address to listen for HTTP requests")
)

func main() {
	flag.Parse()

	fmt.Println("Starting dispatcher")
	model.StartDispatcher(*NFCFS, *NRR, *HRR, *T1, *T2)

	fmt.Println("\nRegistering collector")
	http.HandleFunc("/work", model.Collector)

	fmt.Printf("HTTP Server listening on: %v\n", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}
