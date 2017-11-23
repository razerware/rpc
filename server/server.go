package server

import (
	"net/rpc"
	"log"
	"net"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"strconv"
)

var Datachannel = make(chan MonitorStats, 1)

type MonitorStats struct {
	Hostid     int
	Hostip     string
	CpuPercent float64
	MemPercent float64
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

func (MonitorStats) Collect(stats MonitorStats, res *bool) error {
	Datachannel <- stats
	*res = true
	return nil
}

func ProcessData(datachannel chan MonitorStats) {
	for {
		for c := range datachannel {
			if len(datachannel) >= 0 {
				log.Println("find data")
				go InsertData(c)
			} else {
				log.Println("wait for data")
			}
		}
	}

}

func InsertData(stat MonitorStats) {
	url := "http://127.0.0.1:8086/write?db="
	db := "test"
	url += db
	go sendInfo("cpu", stat, url)
	go sendInfo("mem", stat, url)
}

func sendInfo(field string, stat MonitorStats, url string) {
	measurement := field + "info"
	tags := "hostid=" + strconv.Itoa(stat.Hostid) + ",hostip=" + stat.Hostip
	var stat_string string
	switch field {
	case "cpu":
		stat_string = measurement + "," + tags + " CpuUsage=" + strconv.FormatFloat(stat.CpuPercent, 'f', 2, 64)
	case "mem":
		stat_string = measurement + "," + tags + " MemUsage=" + strconv.FormatFloat(stat.MemPercent, 'f', 2, 64)
	}
	stat_byte := []byte(stat_string)
	fmt.Println(stat_string)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(stat_byte))
	if err != nil {
		// handle error
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	} else {
		log.Println(string(body))
	}
}
