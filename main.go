package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"net/http"
)

var COMMON_PORTS = []string {
	"80",
	"443",
	"1337",
	"3000",
	"4200",
	"7777",
	"8000",
	"8080",
	"8081",
	"8888",
}

const (
	OPEN = 0
	REJECTED = 1
	TIMEOUT = 2
)

type PortResult struct {
	addr string
	status int
}

func scan_port(address string, c chan PortResult) {
	conn, err := net.DialTimeout("tcp", address, time.Duration(2 * time.Second))
	if err != nil {
		log.Printf("WARN: error probing %s %s", address, err)
		c <- PortResult{address, REJECTED}
		return
	}

	defer conn.Close()

	// todo: maybe try to read some http content?
	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	c <- PortResult{address, OPEN}
}

func basic_crawl(w http.ResponseWriter, r *http.Request) {
	c := make(chan PortResult)
	ip := r.RemoteAddr[0: (strings.Index(r.RemoteAddr, ":"))]

	log.Printf("scanning %s", ip)
	for i := 0; i < len(COMMON_PORTS); i++ {
		addr := ip + ":" + COMMON_PORTS[i]
		go scan_port(addr, c)
	}


	for i := 0; i < len(COMMON_PORTS); i++ {
		res := <- c

		if res.status == OPEN {
			fmt.Fprintf(w, "OPEN %s", res.addr)
		}
	}
}

func main() {
	http.HandleFunc("/", basic_crawl)
	log.Fatal(http.ListenAndServe(":80", nil))
}