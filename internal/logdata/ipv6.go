package logdata

import (
	"github.com/google/gopacket/layers"
	"github.com/ma111e/melody/internal/events/helpers"
)

// IPv6LogData is the struct describing the logged data for IPv6 header
type IPv6LogData struct {
	Version        uint8             `json:"version"`
	Length         uint16            `json:"length"`
	NextHeader     layers.IPProtocol `json:"next_header"`
	NextHeaderName string            `json:"next_header_name"`
	TrafficClass   uint8             `json:"traffic_class"`
	FlowLabel      uint32            `json:"flow_label"`
	HopLimit       uint8             `json:"hop_limit"`
	IPLogData      `json:"-"`
}

// NewIPv6LogData is used to create a new IPv6LogData struct
func NewIPv6LogData(ipv6Layer helpers.IPv6Layer) IPv6LogData {
	return IPv6LogData{
		Version:        ipv6Layer.Header.Version,
		Length:         ipv6Layer.Header.Length,
		NextHeader:     ipv6Layer.Header.NextHeader,
		NextHeaderName: ipv6Layer.Header.NextHeader.String(),
		TrafficClass:   ipv6Layer.Header.TrafficClass,
		FlowLabel:      ipv6Layer.Header.FlowLabel,
		HopLimit:       ipv6Layer.Header.HopLimit,
	}
}
