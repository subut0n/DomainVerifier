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
	fmt.Println("Please enter a domain to check MX DNS DMARC DKIM")
	fmt.Println("Example: google.com")
	fmt.Println("------------------------------------------------------------")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC, hasDKIM bool
	var spfRecord, dmarcRecord, dkimDomain string

	// checkMX checks if the given domain has at least one MX record
	hasMX, err := checkMX(domain)
	if err != nil {
		log.Fatalf("Error while checking MX records: %v\n", err)
	}

	// checkSPF checks if the given domain has a valid SPF record
	spfRecord, hasSPF, err = checkSPF(domain)
	if err != nil {
		log.Printf("Error while looking for SPF records: %v\n", err)
	}

	// checkDMARC checks if the given domain has a valid DMARC record
	dmarcRecord, hasDMARC, err = checkDMARC(domain)
	if err != nil {
		log.Printf("Error while looking for DMARC: %v\n", err)
	}

	dkimDomain, hasDKIM, err = checkDKIM(domain)

	// Print results
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("Domain: %s\n", domain)
	fmt.Printf("HasMX: %t\n", hasMX)

	fmt.Printf("HasSPF: %t\n", hasSPF)
	fmt.Printf("SPFRecord: %s\n", spfRecord)

	fmt.Printf("HasDMARC: %t\n", hasDMARC)
	fmt.Printf("DMARCRecord: %s\n", dmarcRecord)

	fmt.Printf("HasDKIM: %t\n", hasDKIM)
	fmt.Printf("DKIMDomain: %s\n", dkimDomain)

	fmt.Println("------------------------------------------------------------")

}

func checkMX(domain string) (bool, error) {
	hasMX := false

	// Lookup MX records for the domain
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return false, err
	}

	// Check if there is at least one MX record
	if len(mxRecords) > 0 {
		hasMX = true
	}

	return hasMX, nil
}

func checkSPF(domain string) (string, bool, error) {
	spfRecord := ""
	hasSPF := false

	// Lookup TXT records for the domain
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return spfRecord, false, err
	}

	// Check if there is a valid SPF record
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	return spfRecord, hasSPF, nil
}

func checkDMARC(domain string) (string, bool, error) {
	dmarcRecord := ""
	hasDMARC := false

	// Lookup TXT records for the DMARC subdomain
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		return dmarcRecord, false, err
	}

	// Check if there is a valid DMARC record
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = domain
			break
		}
	}

	return dmarcRecord, hasDMARC, nil
}

func checkDKIM(domain string) (string, bool, error) {
	dkimRecord := ""
	hasDKIM := false

	// DNS lookup for DKIM record for the given domain
	dkimRecords, err := net.LookupTXT("_domainkey." + domain)
	if err != nil {
		return dkimRecord, false, err
	}

	// Check if the DKIM record is valid
	for _, record := range dkimRecords {
		if strings.HasPrefix(record, "v=DKIM1") {
			hasDKIM = true
			dkimRecord = domain
			break
		}
	}

	return dkimRecord, hasDKIM, nil
}
