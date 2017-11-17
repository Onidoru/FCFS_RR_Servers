package queues

import (
	"fmt"
	"strconv"
	"time"
)

// NewServer creates and returns a new server object
func NewServer(id int, serverQueue chan chan JobRequest, t string, h, t1, t2 float64) Server {
	measure := "ms"
	hD, _ := time.ParseDuration(strconv.FormatFloat(h, 'f', 6, 64) + measure)
	t1D, _ := time.ParseDuration(strconv.FormatFloat(h, 'f', 6, 64) + measure)
	t2D, _ := time.ParseDuration(strconv.FormatFloat(h, 'f', 6, 64) + measure)

	fmt.Println(hD)
	server := Server{
		ID:          id,
		Job:         make(chan JobRequest),
		ServerQueue: serverQueue,
		QuitChan:    make(chan bool),
		Type:        t,
		H:           hD,
		T1:          t1D,
		T2:          t2D,
	}
	return server
}

// Server type
type Server struct {
	ID           int
	Job          chan JobRequest
	FinishedJobs []JobRequest
	DroppedJobs  []JobRequest
	ServerQueue  chan chan JobRequest
	QuitChan     chan bool

	StartTime  time.Time
	FinishTime time.Time

	// Type can be FCFS or RR
	Type string

	// H is the Duration of the process
	// time H, if the job is not finished, it is dropped
	H time.Duration

	// T1 and T2 are relevance measures
	T1, T2 time.Duration
}

// Start starts the server, launching it in the new goroutine
func (s *Server) Start() {
	go func() {
		s.StartTime = time.Now()
		for {
			// add to the server queue
			s.ServerQueue <- s.Job
			select {
			case job := <-s.Job:
				job.TimeInSystem = time.Since(job.ArrivalTime)
				s.reportAboutReceivedJob(job)
				s.process(job)
			case <-s.QuitChan:
				wg.Done()
				wg.Wait()
				s.showStats()
				return
			}
		}
	}()
}

// Stop signalizes to Server to stop receiving jobs
func (s *Server) Stop() {
	go func() {
		s.QuitChan <- true
	}()
}

func (s *Server) process(job JobRequest) {
	switch s.Type {
	case "FCFS":
		time.Sleep(job.Delay)
		job.TimeInProcessing += job.Delay
		job.Relevance = countRelevance(s, job)
		job.TimeInSystem = time.Since(job.ArrivalTime)
		job.IsFinished = true
		s.FinishedJobs = append(s.FinishedJobs, job)
		s.reportAboutFinishedJob(job)
		if len(s.FinishedJobs) == Limit {
			s.Stop()
		}
	case "RR":
		if job.Relevance <= 0 {
			s.DroppedJobs = append(s.DroppedJobs, job)
			s.reportAboutDroppedJob(job)
			if len(s.FinishedJobs)+len(s.DroppedJobs) == Limit {
				s.Stop()
			}
			return
		}
		if job.Delay <= s.H {
			time.Sleep(job.Delay)
			job.TimeInProcessing += job.Delay
			job.Relevance = countRelevance(s, job)
			job.TimeInSystem = time.Since(job.ArrivalTime)
			job.IsFinished = true
			s.FinishedJobs = append(s.FinishedJobs, job)
			s.reportAboutFinishedJob(job)
		} else {
			time.Sleep(s.H)
			job.Delay -= s.H
			job.TimeInProcessing += s.H
			job.Relevance = countRelevance(s, job)
			job.TimeInSystem = time.Since(job.ArrivalTime)
			JobQueueRR <- job
			s.reportAboutDroppedJob(job)
			return
		}
	}
	if job.Id == Limit {
		s.FinishTime = time.Now()
		s.Stop()
	}
}

func (s *Server) reportAboutReceivedJob(j JobRequest) {
	fmt.Println()

	fmt.Printf(`
[%s %v] Received the job:
		- ID: %v
		- Duration: %v
		- Relevance: %v
		- Time In System: %v
		- Time In Processing: %v
		- Arrival Time: %v`,
		s.Type, s.ID,
		j.Id,
		j.Delay,
		j.Relevance,
		j.TimeInSystem,
		j.TimeInProcessing,
		j.ArrivalTime)
}

func (s *Server) reportAboutFinishedJob(j JobRequest) {
	fmt.Println()
	fmt.Printf(`
[%s %v] Finished the job:
		- ID: %v
		- Duration: %v
		- Relevance: %v
		- Time In System: %v
		- Time In Processing: %v
		- Arrival Time: %v`,
		s.Type, s.ID,
		j.Id,
		j.Delay,
		j.Relevance,
		j.TimeInSystem,
		j.TimeInProcessing,
		j.ArrivalTime)
}

func (s *Server) reportAboutDroppedJob(j JobRequest) {
	fmt.Println()
	fmt.Printf(`
[%s %v] Dropped the job:
		- ID: %v
		- Duration: %v
		- Relevance: %v
		- Time In System: %v
		- Time In Processing: %v
		- Arrival Time: %v`,
		s.Type, s.ID,
		j.Id,
		j.Delay,
		j.Relevance,
		j.TimeInSystem,
		j.TimeInProcessing,
		j.ArrivalTime)
}
