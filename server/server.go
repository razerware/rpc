package server

import (
	"net/rpc"
	"log"
	"net"
	"fmt"
	"time"
)

var Datachannel  =make(chan MonitorStats,1)

type MonitorStats struct {
	CpuPercent string
	MemPercent string
}

func NewConn(in interface{}, port int) {
	server := rpc.NewServer()
	log.Printf("Register service:%v\n", in)
	server.Register(in)

	log.Printf("Listen tcp on port %d\n", port)
	address := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Listen error:", err)
	}
	log.Println("Ready to accept connection...")

	conCount := 0
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal("Accept Error:,", err)
				continue
			}
			conCount++
			log.Printf("Receive Client Connection %d\n", conCount)
			go server.ServeConn(conn)
		}
	}()

}

func (MonitorStats)Collect(stats MonitorStats,res *bool)error{
	Datachannel<-stats
	*res=true
	return nil
}

func ProcessData(datachannel chan MonitorStats){
	for  {
		select {
		case c:=<-datachannel:
			fmt.Printf("insert data %v",c)
			InsertData(c)
		}
		time.Sleep(5*time.Second)
	}

}

func InsertData(stat MonitorStats){
	fmt.Printf("insert data",stat)
}