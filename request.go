package main

import (
	"io"
	"log"
	"net"

	"github.com/rcrowley/go-metrics"
)

func NewRequest(in net.Conn, backend, appId string) (err error) {
	var p = Request{backend, appId}
	metrics.GetOrRegisterTimer("request-latency", MetricsRegistry).Time(func() {
		err = p.Accept(in)
	})
	return err
}

type Request struct {
	backend string
	appId   string
}

// Start the request proxy from source -> upstream backend
func (p *Request) Accept(in net.Conn) error {
	defer in.Close()

	out, err := net.Dial("tcp", p.backend)
	defer out.Close()
	if err != nil {
		log.Print("[ERROR] tcp: cannot connect to upstream - ", err)
	}

	// capture all errors in here
	errc := make(chan error, 2)

	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}

	go cp(out, in)
	go cp(in, out)

	err = <-errc
	if err != nil && err != io.EOF {
		log.Print("[WARN]: tcp:  ", err)
		return err
	}
	return nil
}
