package customerimporter

import (
	"testing"
)

func TestFileLoad(t *testing.T) {
	_, err := ScanCSVFileForUniqueEmailDomains("invalid/path.json", ',')
	if err == nil {
		t.Error(err)
	}

	_, err = ScanCSVFileForUniqueEmailDomains("./test.csv", '/')
	if err != nil {
		t.Error(err)
	}
}

func TestCounter(t *testing.T) {
	res, err := ScanCSVFileForUniqueEmailDomains("./test.csv", '/')
	if err != nil {
		t.Error(err)
	}

	if err := res.Error(); err != nil {
		t.Error(err)
	}

	if res.Size != 3 {
		t.Errorf("scanner recognized %d lines (expected 3)", res.Size)
	}
}

func TestOrder(t *testing.T) {
	res, err := ScanCSVFileForUniqueEmailDomains("./test.csv", '/')
	if err != nil {
		t.Error(err)
	}

	if err := res.Error(); err != nil {
		t.Error(err)
	}

	sorted := res.Sort()
	if len(sorted) == 0 {
		t.Error("received empty slice when expected 2 elements")
	}

	if sorted[0].Count <= sorted[1].Count {
		t.Errorf("first element of \"sorted\" slice is smaller than second")
	}
}
