package sys

// Worker interface
type Worker interface {
	// Returns worker name.
	Name() string
	// Starts the worker.
	Start() error
	// Stops the worker.
	Stop()
	// Returns the worker payload.
	Payload() interface{}
	// Returns the worker payload.
	IsReady() bool
	// Enable puts the worker in ready state..
	Enable()
	// Disable puts the worker in not-ready state.
	Disable()
}

// Provider is mail delivery service interface.
type Provider interface {
	Start() error
	// Stop provider.
	Stop() error
	// Client get a provider client implementation.
	Client() interface{}
}
