package wsjtx

import "time"

/*
The heartbeat  message shall be  sent on a periodic  basis every
15   seconds.  This
message is intended to be used by servers to detect the presence
of a  client and also  the unexpected disappearance of  a client
and  by clients  to learn  the schema  negotiated by  the server
after it receives  the initial heartbeat message  from a client.

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l110
*/
type HeartbeatMessage struct {
	Id        string `json:"id"`
	MaxSchema uint32 `json:"maxSchemaVersion"`
	Version   string `json:"version"`
	Revision  string `json:"revision"`
}

/*
WSJT-X  sends this  status message  when various  internal state
changes to allow the server to  track the relevant state of each
client without the need for  polling commands.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l141
*/
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

/*
The decode message is sent when  a new decode is completed, in
this case the 'New' field is true. It is also used in response
to  a "Replay"  message where  each  old decode  in the  "Band
activity" window, that  has not been erased, is  sent in order
as a one of these messages  with the 'New' field set to false.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l206
*/
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

/*
This message is  send when all prior "Decode"  messages in the
"Band Activity"  window have been discarded  and therefore are
no long available for actioning  with a "Reply" message.

The Window  argument  can be  one  of the  following values:

	0  - clear the "Band Activity" window (default)
	1  - clear the "Rx Frequency" window
	2  - clear both "Band Activity" and "Rx Frequency" windows

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l232
*/
type ClearMessage struct {
	Id     string `json:"id"`
	Window uint8  `json:"window"` // In only
}

/*
The QSO logged message is sent when the WSJT-X user accepts the "Log  QSO" dialog by clicking
the "OK" button.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l293
*/
type QsoLoggedMessage struct {
	Id               string    `json:"id"`
	DateTimeOff      time.Time `json:"dateTimeOff"`
	DxCall           string    `json:"dxCall"`
	DxGrid           string    `json:"dxGrid"`
	TxFrequency      uint64    `json:"txFrequency"`
	Mode             string    `json:"mode"`
	ReportSent       string    `json:"reportSent"`
	ReportReceived   string    `json:"reportReceived"`
	TxPower          string    `json:"txPower"`
	Comments         string    `json:"comments"`
	Name             string    `json:"name"`
	DateTimeOn       time.Time `json:"dateTimeOn"`
	OperatorCall     string    `json:"operatorCall"`
	MyCall           string    `json:"myCall"`
	MyGrid           string    `json:"myGrid"`
	ExchangeSent     string    `json:"exchangeSent"`
	ExchangeReceived string    `json:"exchangeReceived"`
}

/*
Close is  sent by  a client immediately  prior to  it shutting
down gracefully.

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l318
*/
type CloseMessage struct {
	Id string `json:"id"`
}

/*
The decode message is sent when  a new decode is completed, in
this case the 'New' field is true.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l381).
*/
type WSPRDecodeMessage struct {
	Id        string  `json:"id"`
	New       bool    `json:"new"`
	Time      uint32  `json:"time"`
	Snr       int32   `json:"snr"`
	DeltaTime float64 `json:"deltaTime"`
	Frequency uint64  `json:"frequency"`
	Drift     int32   `json:"drift"`
	Callsign  string  `json:"callsign"`
	Grid      string  `json:"grid"`
	Power     int32   `json:"power"`
	OffAir    bool    `json:"offAir"`
}

/*
The  logged ADIF  message is  sent to  the server(s)  when the
WSJT-X user accepts the "Log  QSO" dialog by clicking the "OK"
button.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l421
*/
type LoggedAdifMessage struct {
	Id   string `json:"id"`
	Adif string `json:"adif"`
}