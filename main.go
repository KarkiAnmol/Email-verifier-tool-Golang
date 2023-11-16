package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// The main function is the entry point of the program.
func main() {
	// Create a scanner to read input from the console.
	scanner := bufio.NewScanner(os.Stdin)

	// Print the header for the output.
	fmt.Printf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	// Read input from the user until they decide to exit.
	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	// Check for any errors during the scanning process.
	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v \n", err)
	}
}

// The checkDomain function performs domain checks and prints the results.
func checkDomain(domain string) {
	// Declare variables to store the results of domain checks.
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	// Perform a DNS lookup for Mail Exchanger (MX) records.
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}

	// Check if there are any MX records.
	if len(mxRecords) > 0 {
		hasMX = true
	}

	// Perform a DNS lookup for Text (TXT) records.
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}

	// Check each TXT record for the presence of an SPF record.
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	// Perform a DNS lookup for DMARC records.
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
	}

	// Check each DMARC record for the presence of a DMARC record.
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	// Print the results for the current domain.
	fmt.Printf("%v,%v,%v,%v,%v,%v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
