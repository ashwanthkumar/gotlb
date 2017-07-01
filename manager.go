package main

import (
	"log"
	"sync"

	"github.com/ashwanthkumar/golang-utils/maps"
	"github.com/ashwanthkumar/gotlb/providers"
)

// Manager is an abstraction that is responsible for all the frontends
// that TLB manages. Each Frontend represents a particular TCP server
// listening on a specific port and we proxy the requests to one of
// the many backends associated with it
type Manager struct {
	frontends map[string]*Frontend
	lock      sync.Mutex
}

// NewManager returns a new Manager instance which we can Start()
func NewManager() *Manager {
	return &Manager{
		frontends: make(map[string]*Frontend),
	}
}

// Start starts the manager with the given provider
func (m *Manager) Start(provider providers.Provider) {
	addBackend := make(chan providers.BackendInfo)
	removeBackend := make(chan providers.BackendInfo)
	newApp := make(chan providers.AppInfo)
	destroyApp := make(chan providers.AppInfo)
	stopProvider := make(chan bool)

	err := provider.Provide(addBackend, removeBackend, newApp, destroyApp, stopProvider)
	if err != nil {
		log.Fatalf("Unable to start the provider - %v\n", err)
	}

	running := true
	for running {
		select {
		case newBackend := <-addBackend:
			m.AddBackendForApp(newBackend)
		case existingBackend := <-removeBackend:
			m.RemoveBackendForApp(existingBackend)
		case app := <-newApp:
			m.CreateNewFrontendIfNotExist(app)
		case app := <-destroyApp:
			m.RemoveFrontend(app)
		}
	}
}

// RemoveFrontend  removes the specific frontend associated with the app
// it tries to do a graceful shutdown of the frontend
func (m *Manager) RemoveFrontend(app providers.AppInfo) {
	m.lock.Lock()
	defer m.lock.Unlock()
	frontend, present := m.frontends[app.AppId]
	if present {
		frontend.Stop()
		delete(m.frontends, app.AppId)
	}
}

// CreateNewFrontendIfNotExist creates a new frontend and starts it, if one does not exist
// else ignores the app spec associated with it
func (m *Manager) CreateNewFrontendIfNotExist(app providers.AppInfo) {
	m.lock.Lock()
	defer m.lock.Unlock()

	frontend, _ := m.frontends[app.AppId]
	if frontend == nil && maps.Contains(app.Labels, "tlb.port") {
		port := maps.GetString(app.Labels, "tlb.port", "-1")
		frontend = NewFrontend(app.AppId, port, []string{})
		go frontend.Start() // start the frontend
		m.frontends[app.AppId] = frontend
	} else {
		log.Println("[WARN] Either frontend exist else tlb.port does not exist")
	}
}

// AddBackendForApp adds the backend to the list of existing backends for the app
func (m *Manager) AddBackendForApp(backend providers.BackendInfo) {
	frontend, present := m.frontends[backend.AppId]
	if present {
		frontend.AddBackend(backend.Node)
	} else {
		log.Printf("[WARN] Frontend for %s not found. Oops!", backend.AppId)
	}
}

// RemoveBackendForApp removes a specific backend for the app
func (m *Manager) RemoveBackendForApp(backend providers.BackendInfo) {
	frontend, present := m.frontends[backend.AppId]
	if present {
		frontend.RemoveBackend(backend.Node)
	} else {
		log.Printf("[WARN] Frontend for %s not found. Oops!", backend.AppId)
	}
}
