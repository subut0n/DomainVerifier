package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("Domain:\nHasMX:\nHasSPF:\nSPFRecord:\nHasDMARC:\nDMARCRecord:\n")
	fmt.Println("------------------------------------------------------------")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	fmt.Println("------------------------------------------------------------")

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	fmt.Printf("Domain: %s\n", domain)
	fmt.Printf("HasMX: %t\n", hasMX)

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("---\nError: %v\n---\n", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	fmt.Printf("HasSPF: %t\n", hasSPF)
	fmt.Printf("SPFRecord: %s\n", spfRecord)

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		fmt.Println("---")
		log.Printf("Error: %v\n---\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = domain
			break
		}
	}

	fmt.Printf("HasDMARC: %t\n", hasDMARC)
	fmt.Printf("DMARCRecord: %s\n", dmarcRecord)

	fmt.Println("------------------------------------------------------------")
}
