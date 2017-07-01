package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/ashwanthkumar/golang-utils/maps"
	"github.com/ashwanthkumar/gotlb/providers"
	"github.com/ashwanthkumar/gotlb/types"
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
	addBackend := make(chan *types.BackendInfo)
	removeBackend := make(chan *types.BackendInfo)
	newApp := make(chan *types.AppInfo)
	destroyApp := make(chan *types.AppInfo)
	stopProvider := make(chan bool)

	err := provider.Provide(addBackend, removeBackend, newApp, destroyApp, stopProvider)
	if err != nil {
		log.Fatalf("Unable to start the provider - %v\n", err)
	}

	running := true
	for running {
		select {
		case newBackend := <-addBackend:
			err := m.AddBackendForApp(newBackend)
			if err != nil {
				log.Printf("[WARN] %v\n", err)
			}
		case existingBackend := <-removeBackend:
			err := m.RemoveBackendForApp(existingBackend)
			if err != nil {
				log.Printf("[WARN] %v\n", err)
			}
		case app := <-newApp:
			m.CreateNewFrontendIfNotExist(app)
		case app := <-destroyApp:
			m.RemoveFrontend(app)
		}
	}
}

// RemoveFrontend  removes the specific frontend associated with the app
// it tries to do a graceful shutdown of the frontend
func (m *Manager) RemoveFrontend(app *types.AppInfo) {
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
func (m *Manager) CreateNewFrontendIfNotExist(app *types.AppInfo) {
	m.lock.Lock()
	defer m.lock.Unlock()

	frontend, _ := m.frontends[app.AppId]
	if frontend == nil && maps.Contains(app.Labels, types.TLB_PORT) {
		port := maps.GetString(app.Labels, types.TLB_PORT, "-1")
		frontend = NewFrontend(app.AppId, port, []string{})
		go frontend.Start() // start the frontend
		m.frontends[app.AppId] = frontend
	} else {
		log.Println("[WARN] Either frontend exist else tlb.port does not exist")
	}
}

// AddBackendForApp adds the backend to the list of existing backends for the app
func (m *Manager) AddBackendForApp(backend *types.BackendInfo) error {
	frontend, present := m.frontends[backend.AppId]
	if present {
		frontend.AddBackend(backend.Node)
		return nil
	} else {
		return fmt.Errorf("[WARN] Frontend for %s not found. Oops!", backend.AppId)
	}
}

// RemoveBackendForApp removes a specific backend for the app
func (m *Manager) RemoveBackendForApp(backend *types.BackendInfo) error {
	frontend, present := m.frontends[backend.AppId]
	if present {
		frontend.RemoveBackend(backend.Node)
		return nil
	} else {
		return fmt.Errorf("[WARN] Frontend for %s not found. Oops!", backend.AppId)
	}
}

// Used only for tests
func (m *Manager) getFrontend(appId string) (*Frontend, bool) {
	f, exists := m.frontends[appId]
	return f, exists
}

// Used only for tests
func (m *Manager) addFrontend(appId string, frontend *Frontend) {
	m.frontends[appId] = frontend
}
