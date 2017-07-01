package main

import (
	"github.com/ashwanthkumar/golang-utils/sets"
	"github.com/oleiade/lane"
)

// Strategy represents the algorithm that can be
// used to pick a backend to route request to
// General example would be LeastConnection / RoundRobin etc.
type LoadBalancingStrategy interface {
	// Next returns the next backend to route the requests to
	Next() string
	// We need the following 2 methods in order to keep up with
	// the Provider implementation where when a specific backend
	// gets added / removed. Some Strategy implementation requires
	// the to keep the set of backends and some metadata associated
	// with them to return a value in Next()

	// Adds a backend for reference
	AddBackend(backend string)
	// Removes a specific backend for reference
	RemoveBackend(backend string)
}

// LeastConnection is an implementation of Strategy that routes
// requests to a backend based on least number of connections
type LeastConnection struct {
	// TODO - implementat LeastConnection LoadBalancingStrategy
}

// RoundRobin is an implementation of Strategy that routes
// requests to a backend based on round robin fashion
type RoundRobin struct {
	backends        *lane.Queue
	removedBackends sets.Set
}

func RoundRobinStrategy() LoadBalancingStrategy {
	return &RoundRobin{
		backends:        lane.NewQueue(),
		removedBackends: sets.Empty(),
	}
}

func (r *RoundRobin) AddBackend(backend string) {
	r.backends.Enqueue(backend)
}

func (r *RoundRobin) RemoveBackend(backend string) {
	r.removedBackends.Add(backend)
}

func (r *RoundRobin) Next() string {
	item := r.backends.Dequeue().(string)
	if r.removedBackends.Contains(item) {
		// remove the backlist and look again
		r.removedBackends.Remove(item)
		return r.Next()
	} else {
		// add it back at the end of queue so we'll come back to it a little later
		r.backends.Enqueue(item)
		return item
	}
}
