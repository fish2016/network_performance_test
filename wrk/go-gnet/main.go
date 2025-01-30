package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/evanphx/wildcat"
	"github.com/panjf2000/gnet/v2"
)

type httpServer struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
}

type httpCodec struct {
	parser        *wildcat.HTTPParser
	contentLength int
	buf           []byte
}

var (
	CRLF      = []byte("\r\n\r\n")
	lastChunk = []byte("0\r\n\r\n")
)

func (hc *httpCodec) parse(data []byte) (int, []byte, error) {
	bodyOffset, err := hc.parser.Parse(data)
	if err != nil {
		return 0, nil, err
	}

	contentLength := hc.getContentLength()
	if contentLength > -1 {
		bodyEnd := bodyOffset + contentLength
		var body []byte
		if len(data) >= bodyEnd {
			body = data[bodyOffset:bodyEnd]
		}
		return bodyEnd, body, nil
	}

	// Transfer-Encoding: chunked
	if idx := bytes.Index(data[bodyOffset:], lastChunk); idx != -1 {
		bodyEnd := idx + 5
		var body []byte
		if len(data) >= bodyEnd {
			req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(data[:bodyEnd])))
			if err != nil {
				return bodyEnd, nil, err
			}
			body, _ = io.ReadAll(req.Body)
		}
		return bodyEnd, body, nil
	}

	// Requests without a body.
	if idx := bytes.Index(data, CRLF); idx != -1 {
		return idx + 4, nil, nil
	}

	return 0, nil, errors.New("invalid http request")
}

var contentLengthKey = []byte("Content-Length")

func (hc *httpCodec) getContentLength() int {
	if hc.contentLength != -1 {
		return hc.contentLength
	}

	val := hc.parser.FindHeader(contentLengthKey)
	if val != nil {
		i, err := strconv.ParseInt(string(val), 10, 0)
		if err == nil {
			hc.contentLength = int(i)
		}
	}

	return hc.contentLength
}

func (hc *httpCodec) resetParser() {
	hc.contentLength = -1
}

func (hc *httpCodec) reset() {
	hc.resetParser()
	hc.buf = hc.buf[:0]
}

func writeResponse(hc *httpCodec, body []byte) {
	var responseBody []byte
	if string(hc.parser.Path) == "/hello" {
		responseBody = []byte("Hello, World!")
	}

	hc.buf = append(hc.buf, "HTTP/1.1 200 OK\r\nServer: gnet\r\nContent-Type: text/plain\r\nDate: "...)
	hc.buf = time.Now().AppendFormat(hc.buf, "Mon, 02 Jan 2006 15:04:05 GMT")
	hc.buf = append(hc.buf, "\r\nContent-Length: "...)
	hc.buf = append(hc.buf, strconv.Itoa(len(responseBody))...)
	hc.buf = append(hc.buf, "\r\n\r\n"...)
	hc.buf = append(hc.buf, responseBody...)
}

func (hs *httpServer) OnBoot(eng gnet.Engine) gnet.Action {
	hs.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n", hs.multicore, hs.addr)
	return gnet.None
}

func (hs *httpServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(&httpCodec{parser: wildcat.NewHTTPParser()})
	return nil, gnet.None
}

func (hs *httpServer) OnTraffic(c gnet.Conn) gnet.Action {
	hc := c.Context().(*httpCodec)
	buf, _ := c.Peek(-1)
	n := len(buf)

pipeline:
	nextOffset, body, err := hc.parse(buf)
	hc.resetParser()
	if err != nil {
		goto response
	}
	if len(buf) < nextOffset {
		goto response
	}
	writeResponse(hc, body)
	buf = buf[nextOffset:]
	if len(buf) > 0 {
		goto pipeline
	}
response:
	if len(hc.buf) > 0 {
		c.Write(hc.buf)
	}
	hc.reset()
	c.Discard(n - len(buf))
	return gnet.None
}

func main() {
	var port int
	var multicore bool

	// Example command: go run main.go --port 8080 --multicore=true
	flag.IntVar(&port, "port", 8090, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()

	hs := &httpServer{addr: fmt.Sprintf("tcp://127.0.0.1:%d", port), multicore: multicore}

	// Start serving!
	log.Println("server exits:", gnet.Run(hs, hs.addr, gnet.WithMulticore(multicore)))
}
