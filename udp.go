package udp_server

import "net"

type IncomingMessage struct {
	Message []byte
	Remote  *net.UDPAddr
}

func ListenUdp(host string, port int, bufsize int) (chan IncomingMessage, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	incoming := make(chan IncomingMessage)
	buf := make([]byte, bufsize)
	go func() {
		for {
			rlen, remote, err := conn.ReadFromUDP(buf[:])
			if err != nil {
				panic(err)
			}
			incoming <- IncomingMessage{
				buf[0:rlen],
				remote,
			}
		}
	}()
	return incoming, nil
}
