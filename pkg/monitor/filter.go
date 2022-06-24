package monitor

import (
	"github.com/k0swe/wsjtx-go/v4/pkg/callsign"
	"log"
	"strings"
)

type MessageFilter interface {
	Filter(call string) bool
}

type DXCCFilter struct {
	dxcc map[string]int
}

func NewDXCCFilter(dxcc []string) *DXCCFilter {
	dm := make(map[string]int, len(dxcc))
	for _, dxccName := range dxcc {
		dm[strings.TrimSpace(dxccName)] = 0
	}
	return &DXCCFilter{dxcc: dm}
}

func (f DXCCFilter) Filter(message string) bool {
	if de, _, err := callsign.ExtractCallSignFromMessage(message, true); err == nil {
		_, found := f.dxcc[de.DXCC.DXCCName]
		return found
	} else {
		log.Printf("failed to extract callsign %s", err)
	}
	return false
}
