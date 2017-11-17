package queues

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// JobQueues are buffered channels for requests
var JobQueueFCFS = make(chan JobRequest, 3000)
var JobQueueRR = make(chan JobRequest, 3000)
var Id int = 0
var Limit int

// Collector collects jobs to the queue
func Collector(w http.ResponseWriter, r *http.Request) {
	// only POST request is allowed
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if Id == 0 {
		i, _ := strconv.ParseInt(r.FormValue("limit"), 10, 0)
		Limit = int(i)
		Id++
		return
	}

	// parse delay
	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Bad delay value"+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(delay)

	job := JobRequest{
		ArrivalTime: time.Now(),
		Delay:       delay,
		Relevance:   1.0,
		Id:          Id,
		IsFinished:  false,
	}

	Id++

	JobQueueRR <- job
	JobQueueFCFS <- job

	fmt.Println("Job request queued")
	w.WriteHeader(http.StatusCreated)
	return
}
