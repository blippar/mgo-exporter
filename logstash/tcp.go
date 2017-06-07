package logstash

import (
	"encoding/json"
	"net"

	"github.com/apex/log"
)

// TCPForwarder ...
type TCPForwarder struct {
	conn net.Conn
}

// NewTCPForwarder ...
func NewTCPForwarder(connURI string) (*TCPForwarder, error) {

	conn, err := net.Dial("tcp", connURI)
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
	log.WithField("msg", string(js)).Debug("sendMessage")

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
