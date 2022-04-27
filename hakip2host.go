package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// This function grabs the SSL certificate, then dumps the SAN and CommonName
func sslChecks(ip string, resChan chan<- string, client *http.Client) {

	url := ip

	// make sure we use https as we're doing SSL checks
	if strings.HasPrefix(ip, "http://") {
		url = strings.Replace(ip, "http://", "https://", 1)
	} else if !strings.HasPrefix(ip, "https://") {
		url = "https://" + ip
	}

	req, reqErr := http.NewRequest("HEAD", url, nil)
	if reqErr != nil {
		return
	}

	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return
	}

	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		dnsNames := resp.TLS.PeerCertificates[0].DNSNames
		for _, name := range dnsNames {
			resChan <- "[SSL-SAN] " + ip + " " + string(name)
		}
		resChan <- "[SSL-CN] " + ip + " " + resp.TLS.PeerCertificates[0].Subject.CommonName
	}
}

// Do a DNS PTR lookup on the IP
func dnsChecks(ip string, resChan chan<- string, resolver *net.Resolver) {

	addr, err := resolver.LookupAddr(context.Background(), ip)
	if err != nil {
		return
	}

	for _, a := range addr {
		resChan <- "[DNS-PTR] " + ip + " " + a
	}
}

func worker(jobChan <-chan string, resChan chan<- string, wg *sync.WaitGroup, transport *http.Transport, client *http.Client, resolver *net.Resolver) {
	defer wg.Done()

	for job := range jobChan {
		sslChecks(job, resChan, client)
		dnsChecks(job, resChan, resolver)
	}

}
func main() {
	workers := flag.Int("t", 32, "numbers of threads")
	resolverIP := flag.String("r", "", "IP of DNS resolver for lookups")
	dnsProtocol := flag.String("protocol", "udp", "Protocol for DNS lookups (tcp or udp)")
	resolverPort := flag.Int("p", 53, "Port to bother the specified DNS resolver on")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	jobChan := make(chan string)
	resChan := make(chan string)
	done := make(chan struct{})

	// Set up TLS transport
	var transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	// Set up HTTP client
	var client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	// Set up DNS resolver
	var resolver *net.Resolver

	if *resolverIP != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, *dnsProtocol, fmt.Sprintf("%s:%d", *resolverIP, *resolverPort))
			},
		}
	}

	var wg sync.WaitGroup
	wg.Add(*workers)

	go func() {
		wg.Wait()
		close(done)
	}()

	for i := 0; i < *workers; i++ {

		go worker(jobChan, resChan, &wg, transport, client, resolver)
	}

	go func() {
		for scanner.Scan() {
			jobChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
		close(jobChan)
	}()

	for {
		select {
		case <-done:
			return
		case res := <-resChan:
			if strings.HasSuffix(res, ".") {
				res = res[:len(res)-1]
			}
			fmt.Println(res)
		}
	}
}
