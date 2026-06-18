package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

var showVersion bool

var wg sync.WaitGroup

var services = map[int]string{
	20:    "FTP Data",
	21:    "FTP",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	67:    "DHCP Server",
	68:    "DHCP Client",
	69:    "TFTP",
	80:    "HTTP",
	110:   "POP3",
	123:   "NTP",
	135:   "MS RPC",
	137:   "NetBIOS Name Service",
	138:   "NetBIOS Datagram",
	139:   "NetBIOS Session",
	143:   "IMAP",
	161:   "SNMP",
	162:   "SNMP Trap",
	389:   "LDAP",
	443:   "HTTPS",
	445:   "SMB",
	465:   "SMTPS",
	514:   "Syslog",
	587:   "SMTP Submission",
	636:   "LDAPS",
	993:   "IMAPS",
	995:   "POP3S",
	1433:  "Microsoft SQL Server",
	1521:  "Oracle DB",
	1723:  "PPTP",
	2049:  "NFS",
	3306:  "MySQL",
	3389:  "RDP",
	5432:  "PostgreSQL",
	5900:  "VNC",
	6379:  "Redis",
	8080:  "HTTP Proxy / Web",
	8443:  "HTTPS Alt",
	9200:  "Elasticsearch",
	27017: "MongoDB",
}

func scan(host string, port int) {

	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return
	}
	var v string
	if showVersion {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buffer := make([]byte, 1024)
		n, _ := conn.Read(buffer)

		if n > 0 {
			v = (string(buffer[:n]))

		}
	}
	conn.Close()

	if services[port] == "" {
		fmt.Printf("port %d open !service not found! see:https://www.iana.org/assignments/service-names-port-numbers\n", port)
	} else {
		if showVersion {
			fmt.Printf("port %d open with service %s version: %s\n", port, services[port], v)
		}
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
	show_version := flag.Bool("sV", false, "show the version")

	flag.Parse()

	showVersion = *show_version
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
