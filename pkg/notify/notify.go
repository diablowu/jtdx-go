package notify

import (
	"github.com/k0swe/wsjtx-go/v4/pkg/city"
	"log"
)

type Notifier interface {
	Notify(de string, entry *city.DXCCEntry, msg string)
}

type LogPrintNotifier struct {
}

func (n LogPrintNotifier) Notify(de string, entry *city.DXCCEntry, msg string) {
	log.Printf(msg)
}

type WeChatMessageNotifier struct {
}

func (n WeChatMessageNotifier) Notify(de string, entry *city.DXCCEntry, msg string) {

}
