package partialkey

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

var keys [][]byte

func init() {
	for x := 0; x < 1000; x++ {
		keys = append(keys, randStringBytes(8))
	}
	rand.Seed(time.Now().UnixNano())
}

func TestGroup(t *testing.T) {
	g := New(10)

	counters := make([]float64, 10)
	for x := 0; x < len(keys); x++ {
		t := g.Choose(keys[x])
		counters[t]++
	}

	t.Logf("task distribution percentage %0.2f", countersStats(counters, len(keys)))
}

func TestGroupRelease(t *testing.T) {
	g := New(10)

	counters := make([]float64, 10)
	for x := 0; x < len(keys); x++ {
		t := g.Choose(keys[x])
		counters[t]++
		g.Release(t)
	}

	t.Logf("task distribution percentage %0.2f", countersStats(counters, len(keys)))
}

func TestGroupExecTime(t *testing.T) {
	g := New(10)

	var wg sync.WaitGroup
	counters := make([]float64, 10)
	var p []chan int

	for x := 0; x < len(counters); x++ {
		s := make(chan int)
		p = append(p, s)
		go func(s chan int, w int) {
			for range s {
				time.Sleep(time.Duration(w) * time.Millisecond)
				counters[w]++
			}
		}(s, x)
	}

	start := time.Now()
	for x := 0; x < len(keys); x++ {
		wg.Add(1)
		go func(x int) {
			t := g.Choose(keys[x])
			p[t] <- 0
			wg.Done()
		}(x)
	}

	wg.Wait()
	end := time.Since(start)

	t.Logf("Exec time %dms task distribution percentage %0.2f", end.Milliseconds(), countersStats(counters, len(keys)))
}

func TestGroupReleaseExecTime(t *testing.T) {
	g := New(10)

	var wg sync.WaitGroup
	counters := make([]float64, 10)
	var p []chan int

	for x := 0; x < len(counters); x++ {
		s := make(chan int)
		p = append(p, s)
		go func(s chan int, w int) {
			for range s {
				time.Sleep(time.Duration(w) * time.Millisecond)
				counters[w]++
				g.Release(w)
			}
		}(s, x)
	}

	start := time.Now()
	for x := 0; x < len(keys); x++ {
		wg.Add(1)
		go func(x int) {
			t := g.Choose(keys[x])
			p[t] <- 0
			wg.Done()
		}(x)
	}

	wg.Wait()
	end := time.Since(start)

	t.Logf("Exec time %dms task distribution percentage %0.2f", end.Milliseconds(), countersStats(counters, len(keys)))
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	// for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func countersStats(counters []float64, tasks int) []float64 {
	for x := 0; x < len(counters); x++ {
		counters[x] = counters[x] / float64(tasks) * 100
	}

	return counters
}
