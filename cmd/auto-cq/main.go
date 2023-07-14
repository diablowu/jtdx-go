package main

import (
	"bufio"
	"flag"
	"github.com/k0swe/wsjtx-go/v4"
	"github.com/k0swe/wsjtx-go/v4/pkg/monitor"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	bindAddr       = flag.String("bind-addr", "239.255.0.0", "Bind address or Multicast address")
	bindPort       = flag.Uint("bind-port", 2237, "Bind port")
	verbose        = flag.Bool("verbose", false, "Verbose mode")
	autoTxInterval = flag.Duration("auto-tx-interval", time.Minute*5, "Auto trigger tx interval")
)

var defaultDecodeMessageMonitors *monitor.DecodeMessageMonitors

// Simple driver binary for wsjtx-go library.
func main() {
	initCliFlags()

	log.Println("Listening for JTDX...")

	incomingMessageChannel := make(chan interface{}, 5)
	outcomingMessageChannel := make(chan interface{}, 5)

	wsjtxServer, err := wsjtx.MakeServer(*bindAddr, *bindPort)
	if err != nil {
		log.Fatalf("Failed to start auto-cq agent, cause :%v", err)
	}

	errChannel := make(chan error, 5)
	go wsjtxServer.ListenToWsjtx(incomingMessageChannel, outcomingMessageChannel, errChannel)

	ticker := time.NewTicker(*autoTxInterval)
	go func() {
		for {
			<-ticker.C
			log.Println("starting to trigger cq")
			err := wsjtxServer.TriggerCQ(wsjtx.TriggerCQMessage{
				Id:        "JTDX",
				Direction: "",
				TXPeriod:  true,
				Send:      true,
			})
			if err != nil {
				log.Printf("failed to trigger cq %s", err)
			}
		}
	}()

	stdinChannel := make(chan string, 5)
	go stdinCmd(stdinChannel)

	for {
		select {
		case err := <-errChannel:
			log.Printf("error: %v", err)
		case message := <-incomingMessageChannel:
			handleServerMessage(message)
		case command := <-stdinChannel:
			command = strings.ToLower(command)
			handleCommand(command, wsjtxServer)
		}
	}
}

func initCliFlags() {
	flag.Parse()

	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}
}

// Goroutine to listen to stdin.
func stdinCmd(c chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		for scanner.Scan() {
			input := scanner.Text()
			c <- input
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

// When we receive WSJT-X messages, display them.
func handleServerMessage(message interface{}) {
	log.Println("Incoming: ", reflect.TypeOf(message), message)
}

// When we get a command from stdin, send WSJT-X a message.
func handleCommand(command string, wsjtxServer wsjtx.Server) {
	var err error
	switch command {

	case "halt":
		{
			log.Println("Sending Halt")
			err = wsjtxServer.HaltTx(wsjtx.HaltTxMessage{Id: "JTDX"})
		}
	case "hb":
		log.Println("Sending Heartbeat")
		err = wsjtxServer.Heartbeat(wsjtx.HeartbeatMessage{
			Id:        "wsjtx-go",
			MaxSchema: 2,
			Version:   "0.3.1",
			Revision:  "e0d45c929",
		})

	case "clear":
		log.Println("Sending Clear")
		err = wsjtxServer.Clear(wsjtx.ClearMessage{Id: "JTDX"})

	case "cq":
		log.Println("Sending CQ")
		err = wsjtxServer.TriggerCQ(wsjtx.TriggerCQMessage{
			Id:        "JTDX",
			Direction: "",
			TXPeriod:  true,
			Send:      true,
		})

	case "reply":
		log.Println("Sending Replay")
		err = wsjtxServer.Reply(wsjtx.ReplyMessage{
			Id:               "JTDX",
			Time:             uint32(time.Now().Unix()),
			Snr:              -21,
			DeltaTimeSec:     0,
			DeltaFrequencyHz: 50,
			Mode:             "~",
			Message:          "NOCC33 BI1NIZ om90",
			LowConfidence:    false,
		})

	}
	if err != nil {
		log.Println(err)
	}
}
