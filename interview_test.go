// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain. Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
// ================================================================================

package customerimporter

import (
	"fmt"
	"testing"
)

func TestMainCustomersFile(t *testing.T) {
	res, err := ScanCSVFileForUniqueEmailDomains("./customers.csv", ',')
	if err != nil {
		t.Error(err)
	}

	sortedDomainCounters := res.Sort()

	// Do extra data print for showcase
	// I don't have other machines to test it with but on single core (limited to 1400MHz clock speed) it sorted it in 0.002s (without printing to console).
	for _, dc := range sortedDomainCounters {
		fmt.Printf("%s appeared %d time(s)\n", dc.Domain, dc.Count)
	}
}
