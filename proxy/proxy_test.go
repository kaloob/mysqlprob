package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
)

func handlerEchoListener(listener net.Listener) {
	for {

		conn, err := listener.Accept()
		if err != nil {
			break
		}

		go func(conn net.Conn) {
			io.Copy(conn, conn)
			conn.Close()

		}(conn)
	}
}

func TestProxyNoopeInterceptorFn(t *testing.T) {

	source, _ := net.Listen("tcp4", ":40000")
	outer, _ := net.Listen("tcp4", ":40001")
	defer source.Close()
	defer outer.Close()

	go handlerEchoListener(outer)

	fn := func(r io.Reader) ([]byte, error) {
		bytes, _, err := bufio.NewReader(r).ReadLine()

		for i := 0; i < len(bytes); i++ {
			if bytes[i] >= 'a' && bytes[i] <= 'z' {
				bytes[i] -= ('a' - 'A')
			}
		}

		return bytes, err
	}
	proxyHandler := Initialize(source, ":40001", fn)
	go proxyHandler.Run()

	myAddr := new(net.TCPAddr)
	myAddr.IP = net.ParseIP("127.0.0.1")
	myAddr.Port = 50000

	destAddr := new(net.TCPAddr) // It will connects listener of proxyHandler
	destAddr.IP = net.ParseIP("127.0.0.1")
	destAddr.Port = 40000

	conn, err := net.DialTCP("tcp4", myAddr, destAddr)
	if err != nil {
		t.Error(err)
	}

	fmt.Fprintln(conn, "poe")
	conn.CloseWrite()

	bytes, _, err := bufio.NewReader(conn).ReadLine()
	if err != nil {
		t.Error(err)
	}
	conn.CloseRead()

	if !strings.Contains(string(bytes), "POE") {
		t.Errorf("unmatch proxy response. except: %v but got %v ", "POE", string(bytes))
	}

}
