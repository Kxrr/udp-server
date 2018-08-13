package udp_server

import (
	"testing"
	"net"
	"fmt"
	"sync"
	"strings"
)

const port = 9997


type myData struct{
	data []string
	mu sync.Mutex
}


func (r *myData)Add(s string)  {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = append(r.data, s)
}



/**
测试当服务器长时间读取一个连接的消息时, 其它消息是否会丢失
 */
func TestMaxClient(t *testing.T) {
	messages, err := ListenUdp("127.0.0.1", port, 512)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Server started on port %d\n", port)
	originalData := strings.Split("Go is an open source programming language", " ")
	received := myData{}
	sent := myData{}

	go func() {
		for m := range messages {
			received.Add(string(m.Message))
		}
	}()

	wg := &sync.WaitGroup{}
	for i := 0; i < len(originalData); i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
			defer conn.Close()
			if err != nil {
				t.Fatal(err)
			}
			m := fmt.Sprintf(originalData[ii])
			sent.Add(m)
			_, err = conn.Write([]byte(m))
			if err != nil {
				t.Fatal(err)
			}
		}(i)
	}
	wg.Wait()  // all messages sent to server
	fmt.Printf("%#v\n", received.data)
	fmt.Printf("%#v\n", sent.data)
	<-make(chan int, 1)
}
