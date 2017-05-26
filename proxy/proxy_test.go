package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func handlerEchoListener(listener net.Listener) {
	for {

		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			dup_conn := io.TeeReader(conn, os.Stdout)
			io.Copy(conn, dup_conn)
		}(conn)
	}
}

func TestProxyNoopeInterceptorFn(t *testing.T) {

	source, _ := net.Listen("tcp4", ":40000")
	outer, _ := net.Listen("tcp4", ":40001")
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
	handler := Initialize(source, ":40001", fn)
	go handler.Run()

	conn, err := net.DialTimeout("tcp4", ":40000", 5*time.Second)
	defer conn.Close()
	if err != nil {
		t.Error(err)
	}

	fmt.Fprintln(conn, "poe")
	bytes, _, err := bufio.NewReader(conn).ReadLine()
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(bytes), "POE") {
		t.Error("unmatch proxy")
	}
	os.Exit(0)
}
