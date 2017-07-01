package main

import (
	"log"
	"net"
	"sync"
)

// NewFrontend creates a new Frontend instance with appId, frontend
// and array of backends.
func NewFrontend(appId, port string, backends []string) *Frontend {
	return &Frontend{
		appId:    appId,
		backends: backends,
		port:     port,
	}
}

// Frontend represents a instance for an app with a set of backends
type Frontend struct {
	appId    string
	lock     sync.Mutex
	backends []string
	port     string
	listener net.Listener
}

func (f *Frontend) Lookup() string {
	return f.backends[0] // TODO - Replace this with a Strategy implementation for proper load balancing
}

func (f *Frontend) AddBackend(backend string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.backends = append(f.backends, backend)
}

func (f *Frontend) RemoveBackend(backend string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	idx, found := f.findIdxOfBackend(backend)
	if found {
		f.backends = append(f.backends[:idx], f.backends[idx+1:]...)
	} else {
		log.Printf("[WARN] Backend %s is not part of this frontend - %s\n", backend, f.appId)
	}
}

func (f *Frontend) findIdxOfBackend(backend string) (int, bool) {
	for idx, node := range f.backends {
		if node == backend {
			return idx, true
		}
	}

	return -1, false
}

// Start listening on the frontend and start routing requests to backends
func (f *Frontend) Start() {
	log.Printf("Starting Frontend for %s via %s\n", f.appId, f.port)
	l, err := net.Listen("tcp", ":"+f.port)
	f.listener = l
	log.Printf("Started Frontend for %s at %s\n", f.appId, f.port)
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
		go NewRequest(conn, f.Lookup())
	}
}

func (f *Frontend) Stop() {
	log.Println("[INFO] Stopping the frontend - " + f.appId)
	err := f.listener.Close()
	if err != nil {
		log.Printf("[ERR] Error occured while closing the Frontend - %v\n", err)
	}
	log.Println("[INFO] Stopped the frontend - " + f.appId)
}
