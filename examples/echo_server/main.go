package main

import "encoding/binary"
import "github.com/funny/packnet"

// This is an echo server demo work with the echo_client.
// usage:
//     go run github.com/funny/examples/echo_server/main.go
func main() {
	protocol := packnet.NewFixProtocol(4, binary.BigEndian)

	server, err := packnet.ListenAndServe("tcp", "127.0.0.1:10010", protocol)
	if err != nil {
		panic(err)
	}

	server.OnSessionStart(func(session *packnet.Session) {
		println("client", session.RawConn().RemoteAddr().String(), "in")

		session.OnMessage(func(session *packnet.Session, message []byte) {
			println("client", session.RawConn().RemoteAddr().String(), "say:", string(message))

			session.Send(EchoMessage{message})
		})
	})

	server.OnSessionClose(func(session *packnet.Session) {
		println("client", session.RawConn().RemoteAddr().String(), "close")
	})

	server.Start()

	println("server start")

	<-make(chan int)
}

type EchoMessage struct {
	Content []byte
}

func (msg EchoMessage) RecommendPacketSize() uint {
	return uint(len(msg.Content))
}

func (msg EchoMessage) AppendToPacket(packet []byte) []byte {
	return append(packet, msg.Content...)
}
