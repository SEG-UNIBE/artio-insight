package miner

import "sync"

type Queue struct {
	relayMiners []*RelayMiner
	sync.Mutex
}

func (q *Queue) IsEmpty() bool {
	return len(q.relayMiners) == 0
}

func (q *Queue) Length() int {
	return len(q.relayMiners)
}

func (q *Queue) Enqueue(rm *RelayMiner) {
	q.Lock()
	defer q.Unlock()
	q.relayMiners = append(q.relayMiners, rm)
}

func (q *Queue) Dequeue() *RelayMiner {
	q.Lock()
	defer q.Unlock()
	if len(q.relayMiners) == 0 {
		return nil
	}
	rm := q.relayMiners[0]
	q.relayMiners = q.relayMiners[1:]
	return rm
}
