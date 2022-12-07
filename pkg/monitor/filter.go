package monitor

import (
	"github.com/k0swe/wsjtx-go/v4/pkg/callsign"
	"strings"
)

type MessageFilter interface {
	Filter(de *callsign.CallSign) bool
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

func (f DXCCFilter) Filter(de *callsign.CallSign) bool {
	_, found := f.dxcc[de.DXCC.DXCCName]
	return found
}
