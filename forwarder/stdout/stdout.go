package stdout

import (
	"encoding/json"
	"fmt"
)

// Forwarder ...
type Forwarder struct {
	pretty bool
}

// NewForwarder ...
func NewForwarder(prettify bool) (*Forwarder, error) {

	return &Forwarder{
		pretty: prettify,
	}, nil
}

// Send ...
func (f *Forwarder) Send(m interface{}) error {

	var js []byte
	var err error

	if f.pretty {
		if js, err = json.MarshalIndent(m, "", "    "); err != nil {
			return err
		}
	} else {
		if js, err = json.Marshal(m); err != nil {
			return err
		}
	}
	fmt.Println(string(js))

	return nil
}

// Close ...
func (f *Forwarder) Close() error {
	return nil
}
