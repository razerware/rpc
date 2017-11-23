package main

import (
	"rpc/server"
	"os"
	"os/signal"
	"fmt"
)

func main() {
	in := server.MonitorStats{}
	server.NewConn(in, 4200)
	c := make(chan os.Signal)
	signal.Notify(c)

	datachannel:=server.Datachannel
	go server.ProcessData(datachannel)
	for {
		select {
		case <-c:
			fmt.Println("get signal:",c)
			os.Exit(1)
		}
	}
}
