package customerimporter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type SU uint8 // Selected size unit for email repetitions. Switch to uint16, uint32 or uint64 if really needed.

type ScanResult struct {
	Size          uint32
	DomainCounter map[string]SU
	FailedLines   []uint32
}

func (res ScanResult) Error() error {
	if len(res.FailedLines) != 0 {
		s := make([]string, len(res.FailedLines))
		for itx, line := range res.FailedLines {
			s[itx] = strconv.FormatUint(uint64(line), 32)
		}

		return fmt.Errorf("failed to scan \"%s\" line(s)", strings.Join(s, ", "))
	}

	return nil
}

func (res ScanResult) Done() bool {
	return res.Size != 0
}

type DomainCount struct {
	Domain string
	Count  SU
}

func (res *ScanResult) Sort() []DomainCount {
	if res.Size == 0 {
		return make([]DomainCount, 0)
	}

	var size, itx uint32 = uint32(len(res.DomainCounter)), 0
	pairs := make([]DomainCount, size)
	for key, value := range res.DomainCounter {
		pairs[itx] = DomainCount{key, value}
		itx++
	}

	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})

	return pairs
}

// It's a private struct but I want to leave it public for auto documentation reasons.
// However it's "constructor" should be best left private.
func newScanResult() ScanResult {
	return ScanResult{
		Size:          0,
		DomainCounter: make(map[string]SU, 0),
		FailedLines:   make([]uint32, 0),
	}
}
