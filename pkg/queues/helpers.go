package queues

import (
	"fmt"
	"github.com/leesper/go_rng"
	"math/rand"
	// "strconv"
	"math"
	"time"
)

// NormalDistribution returns the pointer to slice with numbers, generated with
// the normal distribution law
// StdDev = 1, Mean = 0
func NormalDistribution(n int, desiredStdDev, desiredMean float64) (sPtr *[]float64) {
	s := []float64{}
	rand.Seed(time.Now().Unix())
	for len(s) < n {
		normalNum := rand.NormFloat64()*desiredStdDev + desiredMean
		if normalNum > 0 {
			s = append(s, normalNum)
		}
	}
	sPtr = &s
	return
}

// PoissonDistribution returns the pointer to slice with numbers, generated
// with the poisson distribution law with parameter l
func PoissonDistribution(lambda float64, n int) (tPtr *[]int64) {
	t := []int64{}
	prng := rng.NewPoissonGenerator(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		t = append(t, prng.Poisson(lambda))
	}
	tPtr = &t
	return
}

func countRelevance(s *Server, j JobRequest) (rel float64) {
	t1 := s.T1.Seconds()
	t2 := s.T2.Seconds()
	tP := j.TimeInProcessing.Seconds()

	if tP <= t1 {
		return 1
	}
	if tP <= t2 {
		rel = (t2 - tP) / (t2 - t1)
		return
	}
	return 0
}

func (s *Server) showStats() {
	x1 := s.countMeanTimeInSystem() * 1000
	x2 := s.countVariance(x1)
	x3 := s.countReactionTime()
	x4 := s.countSuccessRatio()
	x5 := s.countFullRelevance()
	k1, k2, k3, k4, k5 := -1.0, -1.0, -1.0, 3.0, 2.0
	function := k1*x1 + k2*x2 + k3*x3 + k4*x4 + k5*x5
	fmt.Printf(`
[%s %v] Shutdown:
		- Mean Time In System X1: %v
		- Mean Processing Time: %v
		- Variance X2: %v
		- Reaction Time X3: %v
		- Success Ratio X4: %v
		- Full Relevance X5: %v
		- Finale: %v
		- Start Time: %v
		- Finish Time: %v`,
		s.Type, s.ID,
		x1,
		s.countMeanProcessingTime(),
		x2,
		x3,
		x4,
		x5,
		function,
		s.StartTime,
		s.FinishTime)
}

// x1
func (s *Server) countMeanTimeInSystem() float64 {
	var meanTime float64
	for _, job := range s.FinishedJobs {
		meanTime += job.TimeInSystem.Seconds()
	}
	return meanTime / float64(len(s.FinishedJobs))
}

func (s *Server) countUtilization() float64 {
	upTime := s.FinishTime.Sub(s.StartTime).Seconds()
	var workTime float64
	for _, job := range s.FinishedJobs {
		workTime += job.TimeInSystem.Seconds()
	}
	workTime /= float64(len(s.FinishedJobs))
	return workTime / upTime
}

func (s *Server) countMeanProcessingTime() float64 {
	var procTime float64
	for _, job := range s.FinishedJobs {
		if job.IsFinished {
			procTime += job.TimeInProcessing.Seconds()
		}
	}
	return 100.0 * procTime / float64(len(s.FinishedJobs))
}

// x2
func (s *Server) countVariance(x1 float64) float64 {
	var variance float64
	for _, job := range s.FinishedJobs {
		variance += math.Pow((job.TimeInSystem.Seconds()*1000 - x1), 2)
	}
	return variance / float64(len(s.FinishedJobs))
}

// x3
func (s *Server) countReactionTime() float64 {
	var reaction float64
	for _, job := range s.FinishedJobs {
		reaction += s.FinishTime.Sub(s.StartTime).Seconds() - job.TimeInSystem.Seconds()
	}
	return reaction / float64(len(s.FinishedJobs))
}

// x4
func (s *Server) countSuccessRatio() float64 {
	return float64(len(s.FinishedJobs)) / float64(Limit)
}

// x5
func (s *Server) countFullRelevance() float64 {
	var relevance float64
	for _, job := range s.FinishedJobs {
		relevance += job.Relevance
	}
	for _, job := range s.DroppedJobs {
		relevance += job.Relevance
	}
	return relevance
}
