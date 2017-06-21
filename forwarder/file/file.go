package file

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/blippar/mgo-exporter/forwarder"
)

func init() {
	forwarder.Factory.Register(NewForwarder, "file")
}

// Forwarder ...
type Forwarder struct {
	file   *os.File
	pretty bool
}

// NewForwarder ...
func NewForwarder(url *url.URL) (forwarder.Forwarder, error) {

	var prettify bool

	f, err := os.OpenFile(url.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	if url.RawQuery == "pretty" {
		prettify = true
	}

	return &Forwarder{
		file:   f,
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
	js = append(js, byte('\n'))

	f.file.Write(js)

	return nil
}

// Close ...
func (f *Forwarder) Close() error {
	return f.file.Close()
}
