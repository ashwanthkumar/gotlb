package main

import (
	"log"
	"net"
	"sync"
)

// NewFrontend creates a new Frontend instance with appId, frontend
// and array of backends.
func NewFrontend(appId, frontend string, backends []string) *Frontend {
	return &Frontend{
		appId:    appId,
		backends: backends,
		frontend: frontend,
	}
}

// Frontend represents a instance for an app with a set of backends
type Frontend struct {
	appId    string
	lock     sync.Mutex
	backends []string
	frontend string
}

// Start listening on the frontend and start routing requests to backends
func (p *Frontend) Start() {
	log.Printf("Starting Frontend for %s via %s\n", p.appId, p.frontend)
	l, err := net.Listen("tcp", ":"+p.frontend)
	defer l.Close()
	log.Printf("Started Frontend for %s at %s\n", p.appId, p.frontend)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go NewRequest(conn, p.backends[0]) // TODO - Replace this with a Lookup() that does dynamic detection
	}
}
