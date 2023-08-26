package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"strings"


	"golang.org/x/net/idna"
)

var (
	traversedDomains = make(map[string]bool)
)

func main() {
	domainFlag := flag.String("d", "", "Domain to start traversal from")
	depthFlag := flag.Int("n", 1, "Depth of traversal")

	flag.Parse()

	if *domainFlag == "" {
		fmt.Println("Usage: realm -d <domain> -n <depth>")
		return
	}

	domain := *domainFlag
	depth := *depthFlag

	traverseDomain(domain, depth)
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