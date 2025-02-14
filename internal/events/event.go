package events

import (
	"github.com/google/gopacket/layers"
	"github.com/ma111e/melody/internal/loggable"
)

// Event is the interface implementing common methods to generated events
type Event interface {
	//Match(rule rules.Rule) bool
	ToLog() EventLog
	GetIPHeader() *layers.IPv4
	GetICMPv6Header() *layers.ICMPv6
	GetICMPv4Header() *layers.ICMPv4
	GetUDPHeader() *layers.UDP
	GetTCPHeader() *layers.TCP
	GetHTTPData() HTTPEvent

	AddTags(tags map[string]string)
	AddAdditional(add map[string]string)
	loggable.Loggable
}

// EventLog is the interface implementing common methods to generated events' log data
type EventLog interface {
	String() (string, error)
}
