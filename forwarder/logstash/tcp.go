package logstash

import (
	"encoding/json"
	"net"
	"net/url"

	"github.com/blippar/mgo-exporter/forwarder"
)

func init() {
	forwarder.Factory.Register(NewTCPForwarder, "logstash")
}

// TCPForwarder ...
type TCPForwarder struct {
	conn net.Conn
}

// NewTCPForwarder ...
func NewTCPForwarder(connURL *url.URL) (forwarder.Forwarder, error) {

	conn, err := net.Dial("tcp", connURL.Host)
	if err != nil {
		return nil, err
	}

	return &TCPForwarder{
		conn,
	}, nil
}

// Send ...
func (f *TCPForwarder) Send(m interface{}) error {

	var js []byte
	var err error

	if js, err = json.Marshal(m); err != nil {
		return err
	}

	// To work with tls and tcp transports via json_lines codec
	js = append(js, byte('\n'))

	for {
		_, err := f.conn.Write(js)
		if err == nil {
			break
		}
		// if os.Getenv("RETRY_SEND") == "" {
		// 	log.Fatal("logstash: could not write:", err)
		// } else {
		// 	time.Sleep(2 * time.Second)
		// }
		return err
	}
	return nil
}

// Close ...
func (f *TCPForwarder) Close() error {
	return f.conn.Close()
}
