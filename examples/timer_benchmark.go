package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

func runTimers(count int) {
	durationCh := make(chan time.Duration, count)

	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			startedAt := time.Now()
			time.AfterFunc(10*time.Millisecond, func() {
				defer wg.Done()
				durationCh <- time.Since(startedAt)
			})
		}()

	}
	wg.Wait()

	close(durationCh)

	durations := []time.Duration{}
	totalDuration := 0 * time.Millisecond
	for duration := range durationCh {
		durations = append(durations, duration)
		totalDuration += duration
	}
	averageDuration := totalDuration / time.Duration(count)
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	fmt.Printf("run %v timers with average=%v, pct50=%v, pct99=%v\n", count, averageDuration, durations[count/2], durations[int(float64(count)*0.99)])
}

func main() {
	runTimers(10000)
}
