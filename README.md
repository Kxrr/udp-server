# Udp Server

A simple udp server implements in Golang.


## Example


```golang
package main

import (
	"github.com/kxrr/udp-server"
	"fmt"
)

func main()  {
	messages, _ := udp_server.ListenUdp("0.0.0.0", 8881, 1024, 100)
	for m := range messages {
		fmt.Println(string(m.Message))
	}
}
```
