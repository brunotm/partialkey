// Package partialkey is a implementation of the Partial Key Grouping Load-Balanced Partitioning of Distributed Streams algorithm
//
// as described on https://arxiv.org/abs/1510.07623.
// https://melmeric.files.wordpress.com/2014/11/the-power-of-both-choices-practical-load-balancing-for-distributed-stream-processing-engines.pdf
//
// Optionally it includes a small variation to make the balancing take into account the in-flight processing of tasks instead of
// the total allocated in order to prevent contention on slots due to processing time of tasks.
package partialkey

import (
	"sync/atomic"

	"github.com/brunotm/wyfast/go/wyfast"
)

// Group represents a partial key grouping load-balancing algorithm
type Group struct {
	slots []uint64
}

// New creates a new balancing Group with the given slot capacity.
func New(slots int) (g *Group) {
	return &Group{make([]uint64, slots)}
}

// Choose the slot id for the given key and increment it's task count.
// It is safe to call Choose concurrently.
func (k *Group) Choose(key []byte) (slot int) {
	slots := uint64(len(k.slots))
	f := int(wyfast.Sum64(key, 13) % slots)
	s := int(wyfast.Sum64(key, 7) % slots)

	if atomic.LoadUint64(&k.slots[f]) > atomic.LoadUint64(&k.slots[s]) {
		f = s
	}

	atomic.AddUint64(&k.slots[f], 1)
	return f
}

// Release decreases the task count for the given slot.
// Calling release after a task for a given slot makes the key grouping takes into
// account the in-flight requests per slot instead of the total amount of allocated tasks.
// It is safe to call Release concurrently.
func (k *Group) Release(slot int) {
	atomic.AddUint64(&k.slots[slot], ^uint64(0))
}
