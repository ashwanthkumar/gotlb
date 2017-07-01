package main

// Strategy represents the algorithm that can be
// used to pick a backend to route request to
// General example would be LeastConnection / RoundRobin etc.
type Strategy interface {
	// Next returns the next backend to route the requests to
	Next() string
}

// LeastConnection is an implementation of Strategy that routes
// requests to a backend based on least number of connections
type LeastConnection struct{}

// RoundRobin is an implementation of Strategy that routes
// requests to a backend based on round robin fashion
type RoundRobin struct{}
