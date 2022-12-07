package main

import (
	"bufio"
	"flag"
	"github.com/k0swe/wsjtx-go/v4"
	"github.com/k0swe/wsjtx-go/v4/pkg/city"
	"github.com/k0swe/wsjtx-go/v4/pkg/monitor"
	"github.com/k0swe/wsjtx-go/v4/pkg/qywx"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	bindAddr       = flag.String("bind-addr", "239.255.0.0", "Bind address or Multicast address")
	bindPort       = flag.Uint("bind-port", 2237, "Bind port")
	ctyPath        = flag.String("cty-path", "d:/cty.dat", "CTY file")
	verbose        = flag.Bool("verbose", false, "Verbose mode")
	agentID        = flag.Int("agent-id", 1000002, "agent id")
	targetCallSign = flag.String("target-call", "wubo16", "Callsign of received message")
	filteredDXCC   = flag.String("filtered-dxcc", "BY,JA,HL,BV", "Filtered DXCC")
	notifiers      = flag.String("notifiers", "log,wx", "Notifier list")

	autoTxInterval = flag.Duration("auto-tx-interval", time.Minute*5, "Auto trigger tx interval")

	myCall = flag.String("call", "BI1NIZ", "My callsign")
)

var defaultDecodeMessageMonitors *monitor.DecodeMessageMonitors

// Simple driver binary for wsjtx-go library.
func main() {
	flag.Parse()
	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	qywx.Setup(*agentID, *targetCallSign)
	log.Println("Listening for JTDX...")
	if err := city.LoadFromCTYData(*ctyPath); err != nil {
		log.Fatalf("%v", err)
	} else {
		log.Println("Success to load cty data")
	}

	incomingMessageChannel := make(chan interface{}, 5)
	outcomingMessageChannel := make(chan interface{}, 5)
	defaultDecodeMessageMonitors = monitor.CreateDecodeMessageMonitors(
		monitor.NewDefaultMonitor(*myCall, strings.Split(*filteredDXCC, ","), strings.Split(*notifiers, ",")),
		monitor.NewAutoTxTriggerMonitor(*myCall, outcomingMessageChannel))
	//defaultDecodeMessageMonitor =
	//defaultDecodeMessageMonitor =
	wsjtxServer, err := wsjtx.MakeServer(*bindAddr, *bindPort)
	if err != nil {
		log.Fatalf("%v", err)
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
	switch message.(type) {
	//case wsjtx.HeartbeatMessage:
	//	log.Println("Heartbeat:", message)
	case wsjtx.StatusMessage:
		log.Println("Other:", reflect.TypeOf(message), message)
	case wsjtx.DecodeMessage:
		defaultDecodeMessageMonitors.Do(message.(wsjtx.DecodeMessage))
	//case wsjtx.ClearMessage:
	//	log.Println("Clear:", message)
	//case wsjtx.QsoLoggedMessage:
	//	log.Println("QSO Logged:", message)
	//case wsjtx.CloseMessage:
	//	log.Println("Close:", message)
	//case wsjtx.WSPRDecodeMessage:
	//	log.Println("WSPR Decode:", message)
	//case wsjtx.LoggedAdifMessage:
	//	log.Println("Logged Adif:", message)
	default:
		log.Println("Other:", reflect.TypeOf(message), message)
	}
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
			Direction: "DX",
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
