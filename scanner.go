package customerimporter

import (
	"encoding/csv"
	"errors"
	"io"
	"net/mail"
	"os"
	"strings"
)

func ScanCSVFileForUniqueEmailDomains(filePath string, separator rune) (ScanResult, error) {
	if !strings.HasSuffix(strings.ToLower(filePath), ".csv") {
		return ScanResult{}, errors.New("file path points to invalid file type")
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return ScanResult{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 0
	reader.Comma = separator

	// Assume that first record always contains filed names/descriptions.
	// Try to find email field in any of them.
	var emailIndex uint8 = 0
	needFirstIteration := true
	fieldNames, err := reader.Read()
	if err == io.EOF {
		return ScanResult{}, errors.New("targeted csv file is empty")
	}

	if err != nil {
		return ScanResult{}, err
	}

	for itx, fName := range fieldNames {
		if strings.Contains(strings.ToLower(fName), "email") {
			needFirstIteration = false
			emailIndex = uint8(itx)
		}
	}

	res := newScanResult()
	for {
		res.Size++
		records, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return res, err
		}

		if needFirstIteration {
			for itx, fName := range records {
				if _, err := mail.ParseAddress(fName); err == nil {
					needFirstIteration = false
					emailIndex = uint8(itx)
					break
				}
			}

			if needFirstIteration {
				return ScanResult{}, errors.New("failed to detect email field")
			}
		}

		// Check for invalid emails or doubled entry field lines...
		spt := strings.SplitAfter(records[emailIndex], "@")
		if len(spt) < 2 {
			res.FailedLines = append(res.FailedLines, res.Size)
			continue
		}

		res.DomainCounter[spt[1]]++
	}

	res.Size-- // Don't count first line that describes field names/descriptions.
	return res, nil
}
