package main

import (
	"bufio"
	"flag"
	"github.com/k0swe/wsjtx-go/v4/pkg/city"
	"github.com/k0swe/wsjtx-go/v4/pkg/monitor"
	"log"
	"os"
	"strings"

	"github.com/k0swe/wsjtx-go/v4"
)

var (
	jtdxAddr = flag.String("jtdx-addr", "239.255.0.0", "Bind address or Multicast address")
	ctyPath  = flag.String("cty-path", "d:/cty.dat", "CTY file")
)

// Simple driver binary for wsjtx-go library.
func main() {
	flag.Parse()
	//log.SetFlags(log.Llongfile)
	log.Println("Listening for WSJT-X...")
	if err := city.LoadFromCTYData(*ctyPath); err != nil {
		log.Fatalf("%v", err)
	}
	wsjtxServer, err := wsjtx.MakeServer(*jtdxAddr)
	if err != nil {
		log.Fatalf("%v", err)
	}
	wsjtxChannel := make(chan interface{}, 5)
	errChannel := make(chan error, 5)
	go wsjtxServer.ListenToWsjtx(wsjtxChannel, errChannel)

	stdinChannel := make(chan string, 5)
	go stdinCmd(stdinChannel)

	for {
		select {
		case err := <-errChannel:
			log.Printf("error: %v", err)
		case message := <-wsjtxChannel:
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

var defaultDecodeMessageMonitor = monitor.NewDefaultMonitor()

// When we receive WSJT-X messages, display them.
func handleServerMessage(message interface{}) {
	switch message.(type) {
	//case wsjtx.HeartbeatMessage:
	//	log.Println("Heartbeat:", message)
	//case wsjtx.StatusMessage:
	//	log.Println("Status:", message)
	case wsjtx.DecodeMessage:
		defaultDecodeMessageMonitor.Monit(message.(wsjtx.DecodeMessage))
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
		//log.Println("Other:", reflect.TypeOf(message), message)
	}
}

// When we get a command from stdin, send WSJT-X a message.
func handleCommand(command string, wsjtxServer wsjtx.Server) {
	var err error
	switch command {

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
		err = wsjtxServer.Clear(wsjtx.ClearMessage{Id: "WSJT-X", Window: 2})

	case "close":
		log.Println("Sending Close")
		err = wsjtxServer.Close(wsjtx.CloseMessage{Id: "WSJT-X"})

	case "replay":
		log.Println("Sending Replay")
		err = wsjtxServer.Replay(wsjtx.ReplayMessage{Id: "WSJT-X"})

	}
	if err != nil {
		log.Println(err)
	}
}
