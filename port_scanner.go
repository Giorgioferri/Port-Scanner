package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup

var services = map[int]string{
	21:   "FTP",
	22:   "SSH",
	25:   "SMTP",
	53:   "DNS",
	80:   "HTTP",
	443:  "HTTPS",
	3306: "MySQL",
}

func scan(host string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return
	}
	conn.Close()

	if services[port] == "" {
		fmt.Printf("port %d open !service not found! see:https://www.iana.org/assignments/service-names-port-numbers", port)
	} else {
		fmt.Printf("port %d open with service %s\n", port, services[port])
	}
}

func main() {
	host := flag.String("host", "localhost", "host to scan")
	start := flag.Int("start", 1, "start of range") //1 by default
	end := flag.Int("end", 1, "end of range")
	all := flag.Bool("all", false, "scan all 65535 ports")
	flag.Parse()
	if *all {
		*start = 1
		*end = 65535
	}

	for port := *start; port <= *end; port++ {
		wg.Add(1)
		go scan(*host, port, &wg)
	}
	wg.Wait()
	fmt.Println("finished")

}
