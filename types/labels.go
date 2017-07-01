package types

const (
	// Label used to denote the frontend port at which the app is meant to be exposed.
	// This label is mandatory if tlb.enabled = true
	TLB_PORT = "tlb.port"
	// Label used to denote if TCP load balancing is required for this app. Default - false
	TLB_ENABLED = "tlb.enabled"
	// Label used to denote the index of the ports that we should consider while building
	// the backends for the given app. Useful if an app uses multiple ports and want to
	// expose the non-first port via GoTLB. Default - 0
	// This label is a zero-based index.
	TLB_PORTINDEX = "tlb.portIndex"
)
