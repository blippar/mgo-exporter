package forwarder

// Forwarder ...
type Forwarder interface {
	Send(interface{}) error
	Close() error
}
