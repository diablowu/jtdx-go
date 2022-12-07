package monitor

import (
	"fmt"
	"github.com/k0swe/wsjtx-go/v4"
	"github.com/k0swe/wsjtx-go/v4/pkg/callsign"
	"github.com/k0swe/wsjtx-go/v4/pkg/city"
	"github.com/k0swe/wsjtx-go/v4/pkg/notify"
	"log"
	"strings"
)

type DecodeMessageMonitors []DecodeMessageMonitor

func (dms DecodeMessageMonitors) Do(msg wsjtx.DecodeMessage) {
	for _, m := range []DecodeMessageMonitor(dms) {
		m.Monit(msg)
	}
}

func CreateDecodeMessageMonitors(monitors ...DecodeMessageMonitor) *DecodeMessageMonitors {
	var ma []DecodeMessageMonitor
	for _, mo := range monitors {
		ma = append(ma, mo)
	}
	mds := DecodeMessageMonitors(ma)
	return &mds
}

type DecodeMessageMonitor interface {
	Monit(msg wsjtx.DecodeMessage)
}

type dxccFilteredMonitor struct {
	mycall    string
	filter    MessageFilter
	notifiers []notify.Notifier
}

func NewDefaultMonitor(mycall string, filteredDXCC []string, notifiers []string) DecodeMessageMonitor {

	log.Printf("Regiestered notifiers: %v", notify.NotifiersMap)
	log.Printf("Requested notifiers: %v", notifiers)
	var ns []notify.Notifier
	for _, name := range notifiers {
		if n, found := notify.NotifiersMap[strings.TrimSpace(name)]; found {
			ns = append(ns, n())
		}
	}

	log.Printf("Notifiers: %v was be enabled", ns)

	return dxccFilteredMonitor{
		filter:    NewDXCCFilter(filteredDXCC),
		notifiers: ns,
	}
}

func (monitor dxccFilteredMonitor) Monit(msg wsjtx.DecodeMessage) {
	if de, dx, err := callsign.ExtractCallSignFromMessage(msg.Message, true); err == nil {
		deDXCC := city.FindDXCC(de.Number)
		if monitor.filter.Filter(de) {
			if monitor.mycall == dx.Number {
				log.Printf("msg<%s> was call me", msg.Message)
				for _, n := range monitor.notifiers {
					n.Notify(de.Number, deDXCC, fmt.Sprintf("%s", de.Number, msg.Message))
				}
			} else {
				log.Printf("msg<%s> was be filtered", msg.Message)
			}
		} else {
			for _, n := range monitor.notifiers {
				n.Notify(de.Number, deDXCC, fmt.Sprintf("%s(%s),C:%s", de.Number, deDXCC.DXCCName, deDXCC.City))
			}
		}

	} else {
		log.Printf("Failed to ExtractCallSignFromMessage: %s", err)
	}

	//if monitor.filter.Filter(msg.Message) {
	//	log.Printf("msg<%s> was be filtered", msg.Message)
	//} else {
	//	if de, dx, err := callsign.ExtractCallSignFromMessage(msg.Message, true); err == nil {
	//		for _, n := range monitor.notifiers {
	//			dxcc := city.FindDXCC(de.Number)
	//			n.Notify(de.Number, dxcc, fmt.Sprintf("%s(%s),C:%s", de.Number, dxcc.DXCCName, dxcc.City))
	//		}
	//	} else {
	//		log.Printf("Failed to ExtractCallSignFromMessage: %s", err)
	//	}
	//
	//}
}
