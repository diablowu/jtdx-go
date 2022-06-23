package monitor

import (
	"fmt"
	"github.com/k0swe/wsjtx-go/v4"
	"github.com/k0swe/wsjtx-go/v4/pkg/callsign"
	"github.com/k0swe/wsjtx-go/v4/pkg/city"
	"github.com/k0swe/wsjtx-go/v4/pkg/notify"
	"log"
)

type DecodeMessageMonitor interface {
	Monit(msg wsjtx.DecodeMessage)
}

type cjkFilteredMonitor struct {
	filter    MessageFilter
	notifiers []notify.Notifier
}

func NewDefaultMonitor() DecodeMessageMonitor {
	cjkFilter := cjkFilteredMonitor{
		filter:    CJKCallSignFilter{},
		notifiers: make([]notify.Notifier, 2),
	}

	cjkFilter.notifiers[0] = notify.LogPrintNotifier{}
	cjkFilter.notifiers[1] = notify.WeChatMessageNotifier{}
	return cjkFilter
}

func (monitor cjkFilteredMonitor) Monit(msg wsjtx.DecodeMessage) {
	if !monitor.filter.Filter(msg.Message) {
		if de, _, err := callsign.ExtractCallSignFromMessage(msg.Message, true); err == nil {
			for _, n := range monitor.notifiers {
				dxcc := city.FindDXCC(de.Number)
				n.Notify(de.Number, dxcc, fmt.Sprintf("Found DXCC! Call:%s, DXCC:%s, Country:%s", de.Number, dxcc.DXCCName, dxcc.City))
			}
		} else {
			log.Printf("Failed to ExtractCallSignFromMessage: %s", err)
		}
	}
}
