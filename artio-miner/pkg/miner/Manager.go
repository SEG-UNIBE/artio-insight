package miner

import (
	"sync"

	"github.com/SEG-UNIBE/artio-insight/relay-miner/pkg/storage"
)

/*
Manager is the main object to handle the mining process
holds the neo4j instance and the list of miners with some results for for handling recursion
*/
type Manager struct {
	Neo          *storage.Neo4jInstance
	MaxRecursion int
	miners       []*RelayMiner
	loadMap      map[string]bool
	mapMutex     sync.RWMutex
	RelayQueue   *Queue
	MaxRunners   int
	runners      []*Runner
	PushUsers    bool
}

/*
Run starts the mining process for the given relays
the relays are processed in parallel and the information stored in the database
please pay attention to the recursion level to avoid overloading the database with too many requests
*/
func (mgmt *Manager) Run(relays []string) {
	mgmt.loadMap = make(map[string]bool)
	mgmt.mapMutex = sync.RWMutex{}
	mgmt.RelayQueue = new(Queue)

	// push all NIPs
	for i := range 100 {
		mgmt.Neo.Execute(`MERGE(n:NIP {name: $nip})`, map[string]any{"nip": i})
	}

	for _, relay := range relays {
		newMiner := NewMiner(relay)
		newMiner.RecursionLevel = mgmt.MaxRecursion
		mgmt.miners = append(mgmt.miners, newMiner)
	}

	for _, relay := range mgmt.miners {
		mgmt.loadMap[relay.CleanName()] = false
		mgmt.RelayQueue.Enqueue(relay)
	}

	for i := range mgmt.MaxRunners {
		runner := Runner{Manager: mgmt, Id: i}
		mgmt.runners = append(mgmt.runners, &runner)
	}
	mgmt.StartAll()

	for !mgmt.RelayQueue.IsEmpty() || mgmt.AnyNonIdleRunner() {
		// wait until all runners are done
	}
	mgmt.StopAll()
}

func (mgmt *Manager) StartAll() {
	for _, runner := range mgmt.runners {
		go runner.Run()
	}
}

func (mgmt *Manager) StopAll() {
	for _, runner := range mgmt.runners {
		go runner.SignalEnd()
	}
}

func (mgmt *Manager) Dequeue() *RelayMiner {
	return mgmt.RelayQueue.Dequeue()
}

func (mgmt *Manager) Enqueue(rm *RelayMiner) {
	if mgmt.GetLoadMapEntry(rm.CleanName()) {
		// means that we have already processed this relay
		return
	}
	mgmt.SetLoadMapEntryTrue(rm.CleanName())
	mgmt.RelayQueue.Enqueue(rm)
}

/*
AnyRunnerActive checks if any of the provided runners is currently active
*/
func (mgmt *Manager) AnyRunnerActive() bool {
	for _, runner := range mgmt.runners {
		if runner.IsRunning() {
			return true
		}
	}
	return false
}

/*
AnyNonIdleRunner checks if any of the provided runners are all idle or not
*/
func (mgmt *Manager) AnyNonIdleRunner() bool {
	for _, runner := range mgmt.runners {
		if !runner.IsIdle() {
			return true
		}
	}
	return false
}

func (mgmt *Manager) GetLoadMapEntry(relayName string) bool {
	mgmt.mapMutex.RLock()
	defer mgmt.mapMutex.RUnlock()
	return mgmt.loadMap[relayName]
}

func (mgmt *Manager) SetLoadMapEntryTrue(relayName string) {
	mgmt.mapMutex.Lock()
	defer mgmt.mapMutex.Unlock()
	mgmt.loadMap[relayName] = true
}
