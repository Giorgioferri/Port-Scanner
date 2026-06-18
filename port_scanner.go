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
	993:  "IMAPS",
}

func scan(host string, port int) {

	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return
	}
	conn.Close()

	if services[port] == "" {
		fmt.Printf("port %d open !service not found! see:https://www.iana.org/assignments/service-names-port-numbers\n", port)
	} else {
		fmt.Printf("port %d open with service %s\n", port, services[port])
	}
}

func worker(host string, porte chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for porta := range porte {
		scan(host, porta)
	}
}

func main() {
	porte := make(chan int)
	host := flag.String("host", "localhost", "host to scan")
	start := flag.Int("start", 1, "start of range") //1 by default
	end := flag.Int("end", 1, "end of range")
	all := flag.Bool("all", false, "scan all 65535 ports")
	port := flag.Int("port", 0, "scan a single port")

	flag.Parse()
	if *port != 0 {
		*start = *port
		*end = *port
	} else if *all {
		*start = 1
		*end = 65535
	}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(*host, porte, &wg)
	}

	for ports := *start; ports <= *end; ports++ {
		porte <- ports
	}
	close(porte) //dice ai worker basta lavorare
	wg.Wait()
	fmt.Println("finished")

}
