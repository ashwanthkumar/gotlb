package providers

// BackendInfo represents a message from the provider when a new backend is added
// or an existing backend for the app is removed.
type BackendInfo struct {
	AppId string
	Node  string
}

// AppInfo represents the information related to the app
type AppInfo struct {
	AppId  string
	Labels map[string]string
}

// Provider interface defines an implementation that can be used to fetch
// the list of servers for an App. Eg - Marathon, Consul, EtcD, etc.
type Provider interface {
	// Provide gives a set of channels as parameters to the implementation
	// for it to report the respecitve changes accordingly
	// addBackend - Used to denote a particular app instance has been added
	// removeBackend - Used to denote a particular app instance has failed
	// appUpdate - A New app has been deployed / an update to an existing app has been deployed
	// dropApp - An Existing app has been destroyed, we can kill the Frontend for that app
	// stop - Send a value to shutdown the provider, used to gracefully shutdown
	Provide(addBackend chan<- BackendInfo,
		removeBackend chan<- BackendInfo,
		appUpdate chan<- AppInfo,
		dropApp chan<- AppInfo,
		stop <-chan bool) error
}
