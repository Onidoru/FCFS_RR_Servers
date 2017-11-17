package queues

import (
	"fmt"
	"sync"
)

var ServerQueueFCFS chan chan JobRequest
var ServerQueueRR chan chan JobRequest
var wg sync.WaitGroup

// StartDispatcher starts a dispatcher
func StartDispatcher(nFCFS, nRR int, h, t1, t2 float64) {
	ServerQueueFCFS = make(chan chan JobRequest, nFCFS)
	ServerQueueRR = make(chan chan JobRequest, nRR)

	id := 1
	if nFCFS > 0 {
		for ; id <= nFCFS; id++ {
			fmt.Printf("\nStarting FCFS server with ID %v\n", id)
			server := NewServer(id, ServerQueueFCFS, "FCFS", h, t1, t2)
			server.Start()
			wg.Add(1)
		}
	}

	if nRR > 0 {
		for i := id; i < id+nRR; i++ {
			fmt.Printf("\nStarting RR server with ID %v\n", i)
			server := NewServer(i, ServerQueueRR, "RR", h, t1, t2)
			server.Start()
			wg.Add(1)
		}
	}

	go func() {
		for {
			select {
			case job := <-JobQueueFCFS:
				fmt.Println("\n\nFCFS Servers received a job request")
				go func() {
					server := <-ServerQueueFCFS
					fmt.Println("\n\nFCFS System dispatched a job request")
					server <- job
				}()
			case job := <-JobQueueRR:
				fmt.Println("\n\nRR Servers received a job request")
				go func() {
					server := <-ServerQueueRR
					fmt.Println("\n\nRR System dispatched a job request")
					server <- job
				}()
			}
		}
	}()
}
