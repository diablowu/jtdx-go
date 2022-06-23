package monitor

import (
	"github.com/k0swe/wsjtx-go/v4/pkg/callsign"
	"log"
)

type MessageFilter interface {
	Filter(call string) bool
}

type CJKCallSignFilter struct {
}

func (f CJKCallSignFilter) Filter(message string) bool {
	if de, _, err := callsign.ExtractCallSignFromMessage(message, true); err == nil {
		if de.DXCC.DXCCName == "BY" ||
			de.DXCC.DXCCName == "JA" ||
			de.DXCC.DXCCName == "BV" ||
			de.DXCC.DXCCName == "HL" {
			return true
		}
	} else {
		log.Printf("failed to extract callsign %s", err)
	}
	return false
}
