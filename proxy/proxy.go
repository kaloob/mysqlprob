package proxy

import (
	"io"
	"log"
	"net"
	"time"
)

type Proxy struct {
	source      net.Listener
	out         string
	interceptor InterceptorFn
}

type InterceptorFn func(io.Reader) ([]byte, error)

func Initialize(source net.Listener, out string, fn InterceptorFn) *Proxy {
	return &Proxy{source: source, out: out, interceptor: fn}
}

func (p *Proxy) handlerConn(conn net.Conn) {
	outer_conn, err := net.Dial("tcp4", p.out)
	if err != nil {
		panic(err)
	}
	defer outer_conn.Close()
	defer conn.Close()

	for {

		conn.SetReadDeadline(time.Now().Add(time.Second * 1))

		new_bytes, err := p.interceptor(conn)
		if err != nil {
			log.Println(err)
			break

		}
		outer_conn.SetDeadline(time.Now().Add(time.Second * 1))
		write_bytes, err := outer_conn.Write(new_bytes)

		if write_bytes == 0 {
			break
		}

		if err != nil {
			log.Println(err)
			break
		}

		_, err = io.Copy(conn, outer_conn)
		if err != nil {
			log.Println(err)
			break

		}
	}

}

func (p *Proxy) Run() {
	for {
		conn, err := p.source.Accept()

		if err != nil {
			log.Println(err)
			break
		}

		go p.handlerConn(conn)
	}

}
