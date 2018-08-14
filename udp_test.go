package udp_server

import (
	"testing"
	"net"
	"fmt"
	"sync"
	"strings"
	"time"
	"sort"
	"math/rand"
)

var port = 12345 + int(rand.Float32() * 100)

type stringData struct {
	data []string
	mu   sync.Mutex
}

func (r *stringData) Add(s string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = append(r.data, s)
}

func (r *stringData) sortedStrings() string  {
	sort.Sort(sort.StringSlice(r.data))
	return strings.Join(r.data, "|")
}

func (r *stringData) ContainsSameStrings(other stringData) bool {
	if len(r.data) != len(other.data) {
		return false
	}
	return r.sortedStrings() == other.sortedStrings()
}


func TestUDPServer(t *testing.T) {
	messages, err := ListenUdp("127.0.0.1", port, 1024, 100)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("UDP Server started on port %d\n", port)
	originalData := stringData{
		strings.Split(strings.Repeat("Go is an open source programming language ", 100), " "),
		sync.Mutex{},
	}
	received := stringData{}
	go func() {
		for m := range messages {
			received.Add(string(m.Data))

		}
	}()

	sent := stringData{}
	wg := &sync.WaitGroup{}
	for i := 0; i < len(originalData.data); i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
			defer conn.Close()
			if err != nil {
				t.Fatal(err)
			}
			m := fmt.Sprintf(originalData.data[ii])
			sent.Add(m)
			_, err = conn.Write([]byte(m))
			if err != nil {
				t.Fatal(err)
			}
		}(i)
	}
	wg.Wait()                          // waiting the clients sent all messages to server
	time.Sleep(time.Millisecond * 200) // waiting for reading from messages by goroutine
	fmt.Printf("serverReceived = %#v", received.data)
	if !originalData.ContainsSameStrings(sent) {
		t.Fatalf("Expect sent %#v but sent %#v", originalData.sortedStrings(), sent.sortedStrings())
	}
	if !received.ContainsSameStrings(sent) {
		t.Fatalf("Expect server received %#v but received %#v", sent.sortedStrings(), received.sortedStrings())
	}

}
