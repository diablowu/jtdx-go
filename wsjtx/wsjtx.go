package wsjtx

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
)

type HeartbeatMessage struct {
	Id        string `json:"id"`
	MaxSchema uint32 `json:"maxSchemaVersion"`
	Version   string `json:"version"`
	Revision  string `json:"revision"`
}

type StatusMessage struct {
	Id                   string `json:"id"`
	DialFrequency        uint64 `json:"dialFrequency"`
	Mode                 string `json:"mode"`
	DxCall               string `json:"dxCall"`
	Report               string `json:"report"`
	TxMode               string `json:"txMode"`
	TxEnabled            bool   `json:"txEnabled"`
	Transmitting         bool   `json:"transmitting"`
	Decoding             bool   `json:"decoding"`
	RxDF                 uint32 `json:"rxDeltaFreq"`
	TxDF                 uint32 `json:"txDeltaFreq"`
	DeCall               string `json:"deCall"`
	DeGrid               string `json:"deGrid"`
	DxGrid               string `json:"dxGrid"`
	TxWatchdog           bool   `json:"txWatchdog"`
	SubMode              string `json:"submode"`
	FastMode             bool   `json:"fastMode"`
	SpecialOperationMode uint8  `json:"specialMode"`
	FrequencyTolerance   uint32 `json:"frequencyTolerance"`
	TRPeriod             uint32 `json:"txRxPeriod"`
	ConfigurationName    string `json:"configName"`
}

type DecodeMessage struct {
	Id               string  `json:"id"`
	New              bool    `json:"new"`
	Time             uint32  `json:"time"`
	Snr              int32   `json:"snr"`
	DeltaTimeSec     float64 `json:"deltaTime"`
	DeltaFrequencyHz uint32  `json:"deltaFrequency"`
	Mode             string  `json:"mode"`
	Message          string  `json:"message"`
	LowConfidence    bool    `json:"lowConfidence"`
	OffAir           bool    `json:"offAir"`
}

type ClearMessage struct {
	Id string `json:"id"`
}

const Magic = 0xadbccbda

type parser struct {
	buffer []byte
	length int
	cursor int
}

// Goroutine which will listen on a UDP port for messages from WSJT-X. When heard, the messages are
// parsed and then placed in the given channel.
func ListenToWsjtx(c chan interface{}) {
	// TODO: make address and port customizable?
	musticastAddr := "224.0.0.1"
	wsjtxPort := "2237"
	addr, err := net.ResolveUDPAddr("udp", musticastAddr+":"+wsjtxPort)
	check(err)
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	check(err)
	for {
		b := make([]byte, 256)
		length, _, err := conn.ReadFromUDP(b)
		check(err)
		message := ParseMessage(b, length)
		if message != nil {
			c <- message
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Parse messages following the interface laid out in
// https://sourceforge.net/p/wsjt/wsjtx/ci/master/tree/Network/NetworkMessage.hpp. This only parses
// "Out" or "In/Out" message types and does not include "In" types because they will never be
// received by WSJT-X.
func ParseMessage(buffer []byte, length int) interface{} {
	p := parser{buffer: buffer, length: length, cursor: 0}
	magic := p.parseUint32()
	if magic != Magic {
		log.Println("Bad packet")
		return nil
	}
	schema := p.parseUint32()
	if schema != 2 {
		log.Println("Got a schema version I wasn't expecting:", schema)
	}

	messageType := p.parseUint32()
	switch messageType {
	case 0:
		heartbeat := p.parseHeartbeat()
		if p.cursor != p.length {
			log.Printf("Parsing WSJT-X Heartbeat: There were %d bytes left over\n", p.length-p.cursor)
		}
		return heartbeat
	case 1:
		status := p.parseStatus()
		if p.cursor != p.length {
			log.Printf("Parsing WSJT-X Status: There were %d bytes left over\n", p.length-p.cursor)
		}
		return status
	case 2:
		decode := p.parseDecode()
		if p.cursor != p.length {
			log.Printf("Parsing WSJT-X Decode: There were %d bytes left over\n", p.length-p.cursor)
		}
		return decode
	case 3:
		fmt.Print("clear: ")
		fmt.Println(string(p.buffer[p.cursor:]))
	case 5:
		fmt.Print("qso log: ")
		fmt.Println(string(p.buffer[p.cursor:]))
	case 6:
		fmt.Print("close: ")
		fmt.Println(string(p.buffer[p.cursor:]))
	case 10:
		fmt.Print("wspr decode: ")
		fmt.Println(string(p.buffer[p.cursor:]))
	case 12:
		fmt.Print("logged adif: ")
		fmt.Println(string(p.buffer[p.cursor:]))
	}
	return nil
}

func (p *parser) parseHeartbeat() HeartbeatMessage {
	id := p.parseUtf8()
	maxSchema := p.parseUint32()
	version := p.parseUtf8()
	revision := p.parseUtf8()
	return HeartbeatMessage{
		Id:        id,
		MaxSchema: maxSchema,
		Version:   version,
		Revision:  revision,
	}
}

func (p *parser) parseStatus() StatusMessage {
	id := p.parseUtf8()
	dialFreq := p.parseUint64()
	mode := p.parseUtf8()
	dxCall := p.parseUtf8()
	report := p.parseUtf8()
	txMode := p.parseUtf8()
	txEnabled := p.parseBool()
	transmitting := p.parseBool()
	decoding := p.parseBool()
	rxDf := p.parseUint32()
	txDf := p.parseUint32()
	deCall := p.parseUtf8()
	deGrid := p.parseUtf8()
	//TODO: the UDP packets I'm getting don't match the rest of this...
	//dxGrid := p.parseUtf8()
	//txWatchdog := p.parseBool()
	//subMode := p.parseUtf8()
	//fastMode := p.parseBool()
	//specialMode := p.parseUint8()
	//freqTolerance := p.parseUint32()
	//trPeriod := p.parseUint32()
	//configName := p.parseUtf8()
	return StatusMessage{
		Id:            id,
		DialFrequency: dialFreq,
		Mode:          mode,
		DxCall:        dxCall,
		Report:        report,
		TxMode:        txMode,
		TxEnabled:     txEnabled,
		Transmitting:  transmitting,
		Decoding:      decoding,
		RxDF:          rxDf,
		TxDF:          txDf,
		DeCall:        deCall,
		DeGrid:        deGrid,
		//DxGrid:               dxGrid,
		//TxWatchdog:           txWatchdog,
		//SubMode:              subMode,
		//FastMode:             fastMode,
		//SpecialOperationMode: specialMode,
		//FrequencyTolerance:   freqTolerance,
		//TRPeriod:             trPeriod,
		//ConfigurationName:    configName,
	}
}

func (p *parser) parseDecode() DecodeMessage {
	id := p.parseUtf8()
	newDecode := p.parseBool()
	time := p.parseUint32()
	snr := p.parseInt32()
	deltaTime := p.parseFloat64()
	deltaFreq := p.parseUint32()
	mode := p.parseUtf8()
	message := p.parseUtf8()
	lowConfidence := p.parseBool()
	offAir := p.parseBool()
	return DecodeMessage{
		Id:               id,
		New:              newDecode,
		Time:             time,
		Snr:              snr,
		DeltaTimeSec:     deltaTime,
		DeltaFrequencyHz: deltaFreq,
		Mode:             mode,
		Message:          message,
		LowConfidence:    lowConfidence,
		OffAir:           offAir,
	}
}

func (p *parser) parseUint8() uint8 {
	value := p.buffer[p.cursor]
	p.cursor += 1
	return value
}

func (p *parser) parseUtf8() string {
	strlen := int(p.parseUint32())
	value := string(p.buffer[p.cursor:(p.cursor + strlen)])
	p.cursor += strlen
	return value
}

func (p *parser) parseUint32() uint32 {
	value := binary.BigEndian.Uint32(p.buffer[p.cursor : p.cursor+4])
	p.cursor += 4
	return value
}

func (p *parser) parseInt32() int32 {
	value := int32(binary.BigEndian.Uint32(p.buffer[p.cursor : p.cursor+4]))
	p.cursor += 4
	return value
}

func (p *parser) parseUint64() uint64 {
	value := binary.BigEndian.Uint64(p.buffer[p.cursor : p.cursor+8])
	p.cursor += 8
	return value
}

func (p *parser) parseFloat64() float64 {
	bits := binary.BigEndian.Uint64(p.buffer[p.cursor : p.cursor+8])
	value := math.Float64frombits(bits)
	p.cursor += 8
	return value
}

func (p *parser) parseBool() bool {
	value := p.buffer[p.cursor] != 0
	p.cursor += 1
	return value
}