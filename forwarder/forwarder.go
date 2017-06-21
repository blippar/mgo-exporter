package forwarder

import (
	"errors"
	"net/url"
)

// Forwarder ...
type Forwarder interface {
	Send(interface{}) error
	Close() error
}

type forwarderInit func(*url.URL) (Forwarder, error)
type forwarderFactory struct {
	forwarders map[string]forwarderInit
}

// Factory is used to register and create forwarder from an URI
var Factory = forwarderFactory{}

// Register add a forwarder initializer to the factory
func (f *forwarderFactory) Register(fn forwarderInit, name string) {
	if f.forwarders == nil {
		f.forwarders = make(map[string]forwarderInit)
	}
	f.forwarders[name] = fn
}

// NewForwarder create a forwarder from an URI
func (f *forwarderFactory) NewForwarder(uri string) (Forwarder, error) {

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	fw, ok := f.forwarders[u.Scheme]
	if !ok {
		return nil, errors.New("Can't find any forwarder for scheme: " + u.Scheme + "://")
	}
	return fw(u)
}
