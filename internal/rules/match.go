package rules

import (
	"github.com/ma111e/melody/internal/config"
	"github.com/ma111e/melody/internal/events"
)

// Match is the entry point of the Rule matching proc
// It attempt to match every supported rules to the given event after the rules filters have been applied
func (rl *Rule) Match(ev events.Event) bool {
	var noPortsProto = map[string]interface{}{
		config.ICMPv4Kind: struct{}{},
		config.ICMPv6Kind: struct{}{},
	}

	// The rule fails if the source IP is blacklisted
	if len(rl.IPs.BlacklistedIPs) > 0 {
		for _, iprange := range rl.IPs.BlacklistedIPs {
			if iprange.ContainsIPString(ev.GetSourceIP()) {
				return false
			}
		}
	}

	// The rule fails if the source IP is not in the whitelisted addresses
	if len(rl.IPs.WhitelistedIPs) > 0 {
		condOK := false

		for _, iprange := range rl.IPs.WhitelistedIPs {
			if iprange.ContainsIPString(ev.GetSourceIP()) {
				condOK = true
				break
			}
		}

		if !condOK {
			return false
		}
	}

	// If the event's kind supports port filtering
	if _, ok := noPortsProto[ev.GetKind()]; !ok {
		// The rule fails if the source IP is blacklisted
		if len(rl.Ports.BlacklistedPorts) > 0 {
			for _, portRange := range rl.Ports.BlacklistedPorts {
				if portRange.ContainsPort(ev.GetDestPort()) {
					return false
				}
			}
		}

		// The rule fails if the source IP is not in the whitelisted addresses
		if len(rl.Ports.WhitelistedPorts) > 0 {
			condOK := false

			for _, portRange := range rl.Ports.WhitelistedPorts {
				if portRange.ContainsPort(ev.GetDestPort()) {
					condOK = true
					break
				}
			}

			if !condOK {
				return false
			}
		}
	}

	switch ev.GetKind() {
	case config.UDPKind:
		return rl.MatchUDPEvent(ev)
	case config.TCPKind:
		return rl.MatchTCPEvent(ev)
	case config.ICMPv4Kind:
		return rl.MatchICMPv4Event(ev)
	case config.ICMPv6Kind:
		return rl.MatchICMPv6Event(ev)
	case config.HTTPKind:
		fallthrough
	case config.HTTPSKind:
		return rl.MatchHTTPEvent(ev)
	}

	return false
}

// MatchICMPv6Event attempt to match an ICMPv6 event against the calling Rule
func (rl *Rule) MatchICMPv6Event(ev events.Event) bool {
	icmpv6Header := ev.GetICMPv6Header()

	if rl.MatchAll {
		if rl.ICMPv6.Checksum != nil {
			if icmpv6Header.Checksum != *rl.ICMPv6.Checksum {
				return false
			}
		}

		if rl.ICMPv6.TypeCode != nil {
			if uint16(icmpv6Header.TypeCode) != *rl.ICMPv6.TypeCode {
				return false
			}
		}

		if rl.ICMPv6.Code != nil {
			if icmpv6Header.TypeCode.Code() != *rl.ICMPv6.Code {
				return false
			}
		}

		if rl.ICMPv6.Type != nil {
			if icmpv6Header.TypeCode.Type() != *rl.ICMPv6.Type {
				return false
			}
		}

		if rl.ICMPv6.Payload != nil {
			if !rl.ICMPv6.Payload.Match(icmpv6Header.Payload) {
				return false
			}
		}

		return true
	}

	if rl.ICMPv6.Checksum != nil {
		if icmpv6Header.Checksum == *rl.ICMPv6.Checksum {
			return true
		}
	}

	if rl.ICMPv6.TypeCode != nil {
		if uint16(icmpv6Header.TypeCode) == *rl.ICMPv6.TypeCode {
			return true
		}
	}

	if rl.ICMPv6.Code != nil {
		if icmpv6Header.TypeCode.Code() == *rl.ICMPv6.Code {
			return true
		}
	}

	if rl.ICMPv6.Type != nil {
		if icmpv6Header.TypeCode.Type() == *rl.ICMPv6.Type {
			return true
		}
	}

	if rl.ICMPv6.Payload != nil {
		if rl.ICMPv6.Payload.Match(icmpv6Header.Payload) {
			return true
		}
	}

	return false
}

// MatchICMPv4Event attempt to match an ICMPv4 event against the calling Rule
func (rl *Rule) MatchICMPv4Event(ev events.Event) bool {
	icmpv4Header := ev.GetICMPv4Header()

	if rl.MatchAll {
		if rl.ICMPv4.Checksum != nil {
			if icmpv4Header.Checksum != *rl.ICMPv4.Checksum {
				return false
			}
		}

		if rl.ICMPv4.Seq != nil {
			if icmpv4Header.Seq != *rl.ICMPv4.Seq {
				return false
			}
		}

		if rl.ICMPv4.TypeCode != nil {
			if uint16(icmpv4Header.TypeCode) != *rl.ICMPv4.TypeCode {
				return false
			}
		}

		if rl.ICMPv4.Code != nil {
			if icmpv4Header.TypeCode.Code() != *rl.ICMPv4.Code {
				return false
			}
		}

		if rl.ICMPv4.Type != nil {
			if icmpv4Header.TypeCode.Type() != *rl.ICMPv4.Type {
				return false
			}
		}

		if rl.ICMPv4.Payload != nil {
			if !rl.ICMPv4.Payload.Match(icmpv4Header.Payload) {
				return false
			}
		}

		return true
	}

	if rl.ICMPv4.Checksum != nil {
		if icmpv4Header.Checksum == *rl.ICMPv4.Checksum {
			return true
		}
	}

	if rl.ICMPv4.TypeCode != nil {
		if uint16(icmpv4Header.TypeCode) == *rl.ICMPv4.TypeCode {
			return true
		}
	}

	if rl.ICMPv4.Code != nil {
		if icmpv4Header.TypeCode.Code() == *rl.ICMPv4.Code {
			return true
		}
	}

	if rl.ICMPv4.Type != nil {
		if icmpv4Header.TypeCode.Type() == *rl.ICMPv4.Type {
			return true
		}
	}

	if rl.ICMPv4.Seq != nil {
		if icmpv4Header.Seq == *rl.ICMPv4.Seq {
			return true
		}
	}

	if rl.ICMPv4.Payload != nil {
		if rl.ICMPv4.Payload.Match(icmpv4Header.Payload) {
			return true
		}
	}

	return false
}

// MatchUDPEvent attempt to match an UDP event against the calling Rule
func (rl *Rule) MatchUDPEvent(ev events.Event) bool {
	udpHeader := ev.GetUDPHeader()

	if rl.MatchAll {
		//TODO : Add <, > and <> operators
		if rl.UDP.Length != nil {
			if udpHeader.Length != *rl.UDP.Length {
				return false
			}
		}

		if rl.UDP.Checksum != nil {
			if udpHeader.Checksum != *rl.UDP.Checksum {
				return false
			}
		}

		if rl.UDP.Payload != nil {
			if !rl.UDP.Payload.Match(udpHeader.Payload) {
				return false
			}
		}

		//TODO : Add <, > and <> operators
		if rl.UDP.Dsize != nil {
			if uint(len(udpHeader.Payload)) != *rl.UDP.Dsize {
				return false
			}
		}

		return true
	}

	//TODO : Add <, > and <> operators
	if rl.UDP.Length != nil {
		if udpHeader.Length == *rl.UDP.Length {
			return true
		}
	}

	if rl.UDP.Checksum != nil {
		if udpHeader.Checksum == *rl.UDP.Checksum {
			return true
		}
	}

	if rl.UDP.Payload != nil {
		if rl.UDP.Payload.Match(udpHeader.Payload) {
			return true
		}
	}

	//TODO : Add <, > and <> operators
	if rl.UDP.Dsize != nil {
		if uint(len(udpHeader.Payload)) == *rl.UDP.Dsize {
			return true
		}
	}

	return false
}

// MatchTCPEvent attempt to match a TCP event against the calling Rule
func (rl *Rule) MatchTCPEvent(ev events.Event) bool {
	tcpHeader := ev.GetTCPHeader()

	var condOK bool

	if rl.MatchAll {
		if len(rl.TCP.Flags) > 0 {
			condOK = false
			for _, flags := range rl.TCP.Flags {
				// If at least one of the flag string match exactly, then continue
				if tcpHeader.BaseLayer.Contents[13]^(*flags) == 0 {
					condOK = true
					break
				}
			}
			if !condOK {
				return false
			}
		}

		if rl.TCP.Seq != nil {
			if tcpHeader.Seq != *rl.TCP.Seq {
				return false
			}
		}

		if rl.TCP.Ack != nil {
			if tcpHeader.Ack != *rl.TCP.Ack {
				return false
			}
		}

		if rl.TCP.Window != nil {
			if tcpHeader.Window != *rl.TCP.Window {
				return false
			}
		}

		if rl.TCP.Payload != nil {
			if !rl.TCP.Payload.Match(tcpHeader.Payload) {
				return false
			}
		}

		//TODO : Add <, > and <> operators
		if rl.TCP.Dsize != nil {
			if uint(len(tcpHeader.Payload)) != *rl.TCP.Dsize {
				return false
			}
		}

		return true
	}

	// else
	if len(rl.TCP.Flags) > 0 {
		for _, flags := range rl.TCP.Flags {
			// If at least one of the flag string match exactly, then continue
			if tcpHeader.BaseLayer.Contents[13]^(*flags) == 0 {
				return true
			}
		}
	}

	if rl.TCP.Seq != nil {
		if tcpHeader.Seq == *rl.TCP.Seq {
			return true
		}
	}

	if rl.TCP.Ack != nil {
		if tcpHeader.Ack == *rl.TCP.Ack {
			return true
		}
	}

	if rl.TCP.Window != nil {
		if tcpHeader.Window == *rl.TCP.Window {
			return true
		}
	}

	if rl.TCP.Payload != nil {
		if rl.TCP.Payload.Match(tcpHeader.Payload) {
			return true
		}
	}

	//TODO : Add <, > and <> operators
	if rl.TCP.Dsize != nil {
		if uint(len(tcpHeader.Payload)) == *rl.TCP.Dsize {
			return true
		}
	}

	return false
}

// MatchHTTPEvent attempt to match an HTTP event against the calling Rule
func (rl *Rule) MatchHTTPEvent(ev events.Event) bool {
	httpData := ev.GetHTTPData()

	var condOK bool

	if rl.MatchAll {
		if rl.HTTP.URI != nil {
			if !rl.HTTP.URI.Match([]byte(httpData.RequestURI)) {
				return false
			}
		}

		if rl.HTTP.Body != nil {
			if !rl.HTTP.Body.Match([]byte(httpData.Body.Content)) {
				return false
			}
		}

		if rl.HTTP.Headers != nil {
			condOK = false

			for _, inlineHeader := range httpData.InlineHeaders {
				if rl.HTTP.Headers.Match([]byte(inlineHeader)) {
					condOK = true
					break
				}
			}

			if !condOK {
				return false
			}
		}

		if rl.HTTP.Verb != nil {
			if !rl.HTTP.Verb.Match([]byte(httpData.Verb)) {
				return false
			}
		}

		if rl.HTTP.Proto != nil {
			if !rl.HTTP.Proto.Match([]byte(httpData.Proto)) {
				return false
			}
		}

		if rl.HTTP.TLS != nil {
			if *rl.HTTP.TLS != httpData.IsTLS {
				return false
			}
		}
	}

	if rl.HTTP.URI != nil {
		if rl.HTTP.URI.Match([]byte(httpData.RequestURI)) {
			return true
		}
	}

	if rl.HTTP.Body != nil {
		if rl.HTTP.Body.Match([]byte(httpData.Body.Content)) {
			return true
		}
	}

	if rl.HTTP.Headers != nil {
		condOK = false

		for _, inlineHeader := range httpData.InlineHeaders {
			if rl.HTTP.Headers.Match([]byte(inlineHeader)) {
				condOK = true
				break
			}
		}

		if condOK {
			return true
		}
	}

	if rl.HTTP.Verb != nil {
		if rl.HTTP.Verb.Match([]byte(httpData.Verb)) {
			return true
		}
	}

	if rl.HTTP.Proto != nil {
		if rl.HTTP.Proto.Match([]byte(httpData.Proto)) {
			return true
		}
	}

	if rl.HTTP.TLS != nil {
		if *rl.HTTP.TLS == httpData.IsTLS {
			return true
		}
	}

	return false
}
