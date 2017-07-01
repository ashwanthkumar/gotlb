package types

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
