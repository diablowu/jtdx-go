package monitor

import (
	"github.com/k0swe/wsjtx-go/v4"
	"github.com/k0swe/wsjtx-go/v4/pkg/callsign"
	"log"
)

type autoTxTriggerMonitor struct {
	call                    string
	outcomingMessageChannel chan<- interface{}
}

func NewAutoTxTriggerMonitor(call string, outcomingMessageChannel chan<- interface{}) DecodeMessageMonitor {
	return autoTxTriggerMonitor{
		call:                    call,
		outcomingMessageChannel: outcomingMessageChannel,
	}
}

func (monitor autoTxTriggerMonitor) Monit(msg wsjtx.DecodeMessage) {
	log.Printf("Mesage <%s> will be monited.", msg.Message)
	if _, dx, err := callsign.ExtractCallSignFromMessage(msg.Message, true); err == nil {
		if dx.Number == monitor.call {
			monitor.outcomingMessageChannel <- dx.Number
		}
	} else {
		log.Printf("Failed to ExtractCallSignFromMessage: %s", err)
	}
}
