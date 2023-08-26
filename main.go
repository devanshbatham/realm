package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/net/idna"
)

var (
	traversedDomains = make(map[string]bool)
)

func main() {
	domainFlag := flag.String("d", "", "Domain to start traversal from")
	depthFlag := flag.Int("n", 1, "Depth of traversal")
	listFlag := flag.String("l", "", "File containing domain names")

	flag.Parse()

	if *domainFlag == "" && *listFlag == "" {
		fmt.Println("Usage: realm -d <domain> -n <depth> -l <file>")
		return
	}

	if *domainFlag != "" {
		traverseDomain(*domainFlag, *depthFlag)
	}

	if *listFlag != "" {
		domainList, err := readDomainList(*listFlag)
		if err != nil {
			fmt.Println("Error reading domain list:", err)
			return
		}

		for _, domain := range domainList {
			traverseDomain(domain, *depthFlag)
		}
	}
}

func readDomainList(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	var domainList []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			domainList = append(domainList, line)
		}
	}
	return domainList, nil
}

func traverseDomain(domain string, depth int) {
	if traversedDomains[domain] {
		return
	}
	traversedDomains[domain] = true

	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	dnsNames := extractUniqueDNSNames(cert)

	numDNSNames := len(dnsNames)
	yellowPrintf("🔍 Traversing %s: %d domains found\n", domain, numDNSNames)
	printDNSNames(dnsNames)

	if depth > 1 {
		for _, dnsName := range dnsNames {
			traverseDomain(dnsName, depth-1)
		}
	}
}

func printDNSNames(dnsNames []string) {
	for _, dnsName := range dnsNames {
		fmt.Println(dnsName)
	}
}

func extractUniqueDNSNames(cert *x509.Certificate) []string {
	uniqueDNSNames := make(map[string]bool)

	for _, dnsName := range cert.DNSNames {
		asciiDNSName, err := idna.ToASCII(dnsName)
		if err != nil {
			continue
		}

		if strings.HasPrefix(asciiDNSName, "*.") {
			domainWithoutWildcard := strings.TrimPrefix(asciiDNSName, "*.")
			uniqueDNSNames[domainWithoutWildcard] = true
		} else {
			uniqueDNSNames[asciiDNSName] = true
		}
	}

	var dnsNames []string
	for dnsName := range uniqueDNSNames {
		dnsNames = append(dnsNames, dnsName)
	}

	return dnsNames
}

func yellowPrintf(format string, a ...interface{}) {
	yellow := "\033[33m%s\033[0m"
	fmt.Printf(yellow, fmt.Sprintf(format, a...))
}
