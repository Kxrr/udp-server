package udp_server

import (
	"net"
	"errors"
)

type IncomingMessage struct {
	Data []byte
	Remote  *net.UDPAddr
	Error error
}


var ErrorOverFlow = errors.New("data overflowed, maybe you can increase the dataBuf")

func ListenUdp(host string, port int, dataBuf int, incomingBuf int) (chan IncomingMessage, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	// use a buffered channel here
	// data from connection will queued if no one read from the channel
	incoming := make(chan IncomingMessage, incomingBuf)
	buf := make([]byte, dataBuf)
	go func() {
		for {
			rlen, remote, err := conn.ReadFromUDP(buf)
			if err != nil {
				incoming <- IncomingMessage{
					nil,
					remote,
					err,
				}
				continue
			}
			data := make([]byte, rlen)
			copy(data, buf[0:rlen])
			if rlen > dataBuf {
				incoming <- IncomingMessage{
					data,
					remote,
					ErrorOverFlow,
				}
				continue
			}
			incoming <- IncomingMessage{
				data,
				remote,
				nil,
			}
		}
	}()
	return incoming, nil
}
