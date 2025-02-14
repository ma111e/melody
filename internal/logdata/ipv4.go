package logdata

import (
	"strings"

	"github.com/google/gopacket/layers"
	"github.com/ma111e/melody/internal/events/helpers"
)

// IPv4LogData is the struct describing the logged data for IPv4 header
type IPv4LogData struct {
	Version    uint8             `json:"version"`
	IHL        uint8             `json:"ihl"`
	TOS        uint8             `json:"tos"`
	Length     uint16            `json:"length"`
	ID         uint16            `json:"id"`
	Fragbits   string            `json:"fragbits"`
	FragOffset uint16            `json:"frag_offset"`
	TTL        uint8             `json:"ttl"`
	Protocol   layers.IPProtocol `json:"protocol"`
	IPLogData  `json:"-"`
}

// NewIPv4LogData is used to create a new IPv4LogData struct
func NewIPv4LogData(ipv4Layer helpers.IPv4Layer) IPv4LogData {
	var ipFlagsStr []string

	if ipv4Layer.Header.Flags&layers.IPv4EvilBit != 0 {
		ipFlagsStr = append(ipFlagsStr, "EV")
	}
	if ipv4Layer.Header.Flags&layers.IPv4DontFragment != 0 {
		ipFlagsStr = append(ipFlagsStr, "DF")
	}
	if ipv4Layer.Header.Flags&layers.IPv4MoreFragments != 0 {
		ipFlagsStr = append(ipFlagsStr, "MF")
	}

	data := IPv4LogData{
		Version:    ipv4Layer.Header.Version,
		IHL:        ipv4Layer.Header.IHL,
		TOS:        ipv4Layer.Header.TOS,
		Length:     ipv4Layer.Header.Length,
		ID:         ipv4Layer.Header.Id,
		FragOffset: ipv4Layer.Header.FragOffset,
		TTL:        ipv4Layer.Header.TTL,
		Protocol:   ipv4Layer.Header.Protocol,
		Fragbits:   strings.Join(ipFlagsStr, ""),
	}

	return data
}
