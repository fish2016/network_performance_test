package main

import (
	"bytes"
	"os"
	"sync"
	"time"

	uringnet "github.com/y001j/uringnet"
	socket "github.com/y001j/uringnet/sockets"
)

type testServer struct {
	uringnet.BuiltinEventEngine

	testloop *uringnet.Ringloop
	//ring      *uring_net.URingNet
	addr      string
	multicore bool
}

type httpCodec struct {
	delimiter []byte
	buf       []byte
}

func appendResponse(hc *[]byte) {
	*hc = append(*hc, "HTTP/1.1 200 OK\r\nServer: uringNet\r\nContent-Type: text/plain\r\nDate: "...)
	*hc = time.Now().AppendFormat(*hc, "Mon, 02 Jan 2006 15:04:05 GMT")
	*hc = append(*hc, "\r\nContent-Length: 12\r\n\r\nHello World!"...)
}

var (
	errMsg      = "Internal Server Error"
	errMsgBytes = []byte(errMsg)
)

func (ts *testServer) OnTraffic(data *uringnet.UserData, ringnet *uringnet.URingNet) uringnet.Action {

	//将data.Buffer转换为string
	//buffer := data.Buffer[:data.BufSize]

	buffer := ringnet.ReadBuffer
	//tes :=
	//fmt.Println("data:", " offset: ", tes, " ", data.BufOffset)
	//获取tes中“\r\n\r\n”的数量
	count := bytes.Count(buffer, []byte("GET"))
	if count == 0 {
		//appendResponse(&data.WriteBuf)
		//return UringNet.Close
		return uringnet.None
	} else {
		for i := 0; i < count; i++ {
			appendResponse(&data.WriteBuf)
		}
	}
	return uringnet.Echo
}

func (ts *testServer) OnWritten(data uringnet.UserData) uringnet.Action {

	return uringnet.None
}

func (ts *testServer) OnOpen(data *uringnet.UserData) ([]byte, uringnet.Action) {

	ts.SetContext(&httpCodec{delimiter: []byte("\r\n\r\n")})
	return nil, uringnet.None
}

func main() {
	addr := os.Args[1]
	//runtime.GOMAXPROCS(runtime.NumCPU()*2 - 1)

	options := socket.SocketOptions{TCPNoDelay: socket.TCPNoDelay, ReusePort: true}
	ringNets, _ := uringnet.NewMany(uringnet.NetAddress{socket.Tcp4, addr}, 3200, true, 3, options, &testServer{}) //runtime.NumCPU()

	loop := uringnet.SetLoops(ringNets, 4000)

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)

	loop.RunMany()

	waitgroup.Wait()
}
