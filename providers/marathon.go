package providers

import (
	"fmt"
	"log"
	"net/url"

	"github.com/ashwanthkumar/golang-utils/maps"
	"github.com/ashwanthkumar/golang-utils/sets"
	marathon "github.com/gambol99/go-marathon"
)

type MarathonProvider struct {
	addBackend    chan<- BackendInfo
	removeBackend chan<- BackendInfo
	appUpdate     chan<- AppInfo
	dropApp       chan<- AppInfo
	stopMe        <-chan bool
	apps          sets.Set

	marathonHost string
}

func NewMarathonProvider(marathonHost string) Provider {
	return &MarathonProvider{
		marathonHost: marathonHost,
		apps:         sets.Empty(),
	}
}

func (m *MarathonProvider) Provide(
	addBackend chan<- BackendInfo,
	removeBackend chan<- BackendInfo,
	appUpdate chan<- AppInfo,
	dropApp chan<- AppInfo,
	stop <-chan bool) error {
	m.addBackend = addBackend
	m.removeBackend = removeBackend
	m.appUpdate = appUpdate
	m.dropApp = dropApp
	m.stopMe = stop
	log.Println("Starting Marathon Provider on " + m.marathonHost)
	go m.start()
	log.Println("Marathon Provider Started and configured to " + m.marathonHost)
	return nil
}

func (m *MarathonProvider) start() {
	config := marathon.NewDefaultConfig()
	config.URL = m.marathonHost
	config.EventsTransport = marathon.EventsTransportSSE
	client, err := marathon.NewClient(config)
	if err != nil {
		log.Fatalf("Unable to create marathon client - %v\n", err)
	}

	// Initialize all the apps since we're bootstrapping
	m.lookOverAllApps(client)

	eventsChannel, err := client.AddEventsListener(marathon.EventIDAPIRequest | marathon.EventIDStatusUpdate | marathon.EventIDFailedHealthCheck | marathon.EventIDAppTerminated)
	if err != nil {
		log.Fatalf("Unable to create events listener - %v\n", err)
	}

	running := true
	for running {
		select {
		case event := <-eventsChannel:
			switch event.ID {
			case marathon.EventIDStatusUpdate:
				update := event.Event.(*marathon.EventStatusUpdate)
				// check if the update is for known app
				knownApp := m.apps.Contains(update.AppID)

				if knownApp && update.TaskStatus == "TASK_FAILED" {
					m.removeBackend <- BackendInfo{
						AppId: update.AppID,
						// TODO - Support choosing different ports / ip address
						Node: update.IPAddresses[0].IPAddress + ":" + fmt.Sprintf("%d", update.Ports[0]),
					}
				} else if knownApp && update.TaskStatus == "TASK_RUNNING" {
					m.addBackend <- BackendInfo{
						AppId: update.AppID,
						// TODO - Support choosing different ports / ip address
						Node: update.IPAddresses[0].IPAddress + ":" + fmt.Sprintf("%d", update.Ports[0]),
					}
				}
				// fmt.Printf("app=%s, id=%s, slaveId=%s, status=%s, host:ip=%s:%d\n", update.AppID, update.TaskID, update.SlaveID, update.TaskStatus, update.IPAddresses[0].IPAddress, update.Ports[0])
			case marathon.EventIDAPIRequest:
				appRequest := event.Event.(*marathon.EventAPIRequest)
				_, err := client.Application(appRequest.AppDefinition.ID)
				if err != nil {
					log.Printf("[WARN] Unable to get application - %s - %v\n", appRequest.AppDefinition.ID, err)
					fmt.Printf("Deleted the App spec - %v\n", appRequest)
					// check if the update is for known app
					knownApp := m.apps.Contains(appRequest.AppDefinition.ID)
					if knownApp {
						// most likely the app was destroyed
						m.dropApp <- AppInfo{
							AppId:  appRequest.AppDefinition.ID,
							Labels: *appRequest.AppDefinition.Labels,
						}
					}
				} else {
					fmt.Printf("New / Updated the App spec - %v\n", appRequest)
					m.appUpdate <- AppInfo{
						AppId:  appRequest.AppDefinition.ID,
						Labels: *appRequest.AppDefinition.Labels,
					}
				}
			}
		case <-m.stopMe:
			running = false
			client.RemoveEventsListener(eventsChannel)
		}
	}
}

func (m *MarathonProvider) lookOverAllApps(client marathon.Marathon) {
	v := url.Values{}
	v.Set("embed", "apps.tasks")
	apps, err := client.Applications(v)
	if err != nil {
		log.Printf("[WARN] Initializing with all applications failed - %v\n", err)
	} else {
		for _, app := range apps.Apps {
			if maps.GetBoolean(*app.Labels, "tlb.enabled", false) {
				log.Printf("Adding new app - %s\n", app.ID)
				m.appUpdate <- AppInfo{
					AppId:  app.ID,
					Labels: *app.Labels,
				}
				// add this app to the list of known apps
				m.apps.Add(app.ID)
				portIndex := maps.GetInt(*app.Labels, "tlb.portIndex", 0)
				log.Printf("[DEBUG] portIndex used for %s is %d\n", app.ID, portIndex)
				// add the list of tasks to update the backend
				for _, task := range app.Tasks {
					log.Printf("[DEBUG] Adding backend for %s as %v\n", app.ID, task)
					m.addBackend <- BackendInfo{
						AppId: app.ID,
						Node:  task.IPAddresses[portIndex].IPAddress + ":" + fmt.Sprintf("%d", task.Ports[portIndex]),
					}
				}
			}
		}
	}
}
