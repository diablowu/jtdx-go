package main

import (
	"flag"
	"github.com/k0swe/wsjtx-go/v4"
	"github.com/k0swe/wsjtx-go/v4/pkg/monitor"
	"log"
	"time"
)

var (
	jtdxPort       = flag.Uint("jtdx-port", 2237, "Bind port")
	duration       = flag.Duration("auto-tx-duration", time.Minute*50, "Auto trigger tx duration")
	autoTxInterval = flag.Duration("auto-tx-interval", time.Minute*5, "Auto trigger tx interval")
)

var defaultDecodeMessageMonitor monitor.DecodeMessageMonitor

func main() {
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	ticker := time.NewTicker(*autoTxInterval)
	durationTimer := time.NewTimer(*duration)

	jtdxClient, err := wsjtx.MakeClient(*jtdxPort)

	if err != nil {
		panic(err)
	}

	running := true

	for running {
		select {
		case <-ticker.C:
			{
				log.Println("starting to cq")
				err := jtdxClient.TriggerCQ(wsjtx.TriggerCQMessage{
					Id:        "JTDX",
					Direction: "",
					TXPeriod:  true,
					Send:      true,
				})

				if err != nil {
					log.Println(err)
				}
			}
		case <-durationTimer.C:
			{
				running = false
				break
			}
		}
	}

	log.Println("Auto TX was finished.")
}
