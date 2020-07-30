package partialkey

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestGroup(t *testing.T) {
	g := New(10)
	var keys [][]byte
	for x := 0; x < 1000; x++ {
		keys = append(keys, randStringBytes(8))
	}

	counters := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for x := 0; x < 1000; x++ {
		counters[g.Choose(keys[x])]++
	}

	for x := 0; x < 10; x++ {
		counters[x] = counters[x] / 1000 * 100
	}

	t.Logf("task distribution percentage %0.2f", counters)
}

func TestGroupRelease(t *testing.T) {
	g := New(10)
	var keys [][]byte
	for x := 0; x < 1000; x++ {
		keys = append(keys, randStringBytes(8))
	}

	counters := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for x := 0; x < 1000; x++ {
		t := g.Choose(keys[x])
		counters[t]++
		g.Release(t)
	}

	for x := 0; x < 10; x++ {
		counters[x] = counters[x] / 1000 * 100
	}

	t.Logf("task distribution percentage %0.2f", counters)
}

func TestGroupExecTime(t *testing.T) {
	g := New(10)
	var keys [][]byte
	for x := 0; x < 1000; x++ {
		keys = append(keys, randStringBytes(8))
	}

	var wg sync.WaitGroup
	counters := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for x := 0; x < 1000; x++ {
		wg.Add(1)
		t := g.Choose(keys[x])
		counters[t]++
		go func() {
			<-time.After(time.Millisecond * 2 * time.Duration(rand.Intn(1000)))
			wg.Done()
		}()

	}

	wg.Wait()

	for x := 0; x < 10; x++ {
		counters[x] = counters[x] / 1000 * 100
	}

	t.Logf("task distribution percentage %0.2f", counters)
}

func TestGroupReleaseExecTime(t *testing.T) {
	g := New(10)
	var keys [][]byte
	for x := 0; x < 1000; x++ {
		keys = append(keys, randStringBytes(8))
	}

	var wg sync.WaitGroup
	counters := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for x := 0; x < 1000; x++ {
		wg.Add(1)
		t := g.Choose(keys[x])
		counters[t]++
		go func() {
			<-time.After(time.Millisecond * 2 * time.Duration(rand.Intn(1000)))
			g.Release(t)
			wg.Done()
		}()

	}

	wg.Wait()

	for x := 0; x < 10; x++ {
		counters[x] = counters[x] / 1000 * 100
	}

	t.Logf("task distribution percentage %0.2f", counters)
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
