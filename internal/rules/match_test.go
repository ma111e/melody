package rules

import (
	"bufio"
	"bytes"
	"net/http"
	"testing"

	"github.com/ma111e/melody/internal/config"

	"github.com/google/gopacket/layers"
	"github.com/ma111e/melody/internal/events"
)

func init() {
	config.Cfg = config.NewConfig()
}

func TestMatchingLogicFlow(t *testing.T) {
	ruleFilename := "logic_flow_rules.yml"
	pcapFilename := "logic_flow.pcap"
	rawPackets := false
	var rule Rule

	ruleset, err := LoadRuleFile(ruleFilename)
	if err != nil {
		t.Error(err)
		return
	}

	filteredEvents, _, err := ReadPacketsFromPcap(pcapFilename, layers.IPProtocolTCP, rawPackets)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filteredEvents) == 0 {
		t.Error("No TCP packets has been read from pcap")
		return
	}

	tests := []struct {
		Ok     []string
		Nok    []string
		Packet events.Event
	}{
		{
			Ok: []string{
				"ok_any_sub",
				"ok_all_sub",
				"ok_all_upper",
				"ok_any_upper",
				"ok_any_upper_mixed",
				"ok_all_upper_mixed",
				"ok_all_all_full_mixed",
				"ok_all_any_full_mixed",
				"ok_any_any_full_mixed",
			},
			Nok: []string{
				"nok_any_sub",
				"nok_all_sub",
				"nok_all_upper",
				"nok_any_upper",
				"nok_any_upper_mixed",
				"nok_all_upper_mixed",
				"nok_all_all_full_mixed",
				"nok_any_any_full_mixed",
			},
			Packet: filteredEvents[2],
		},
	}

	for _, suite := range tests {
		for _, rulename := range suite.Ok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); !ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
		for _, rulename := range suite.Nok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
	}
}

func TestMatchUDPEvent(t *testing.T) {
	ruleFilename := "udp_rules.yml"
	pcapFilename := "udp_values.pcap"
	rawPackets := false
	var rule Rule

	ruleset, err := LoadRuleFile(ruleFilename)
	if err != nil {
		t.Error(err)
		return
	}

	filteredEvents, _, err := ReadPacketsFromPcap(pcapFilename, layers.IPProtocolUDP, rawPackets)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filteredEvents) == 0 {
		t.Error("No UDP packets has been read from pcap")
		return
	}

	tests := []struct {
		Ok     []string
		Nok    []string
		Packet events.Event
	}{
		{
			Ok: []string{
				"ok_checksum",
				"ok_dsize",
				"ok_length",
				"ok_payload",
			},
			Nok: []string{
				"nok_checksum",
				"nok_dsize",
				"nok_length",
				"nok_payload",
			},
			Packet: filteredEvents[0],
		},
	}

	for _, suite := range tests {
		for _, rulename := range suite.Ok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); !ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
		for _, rulename := range suite.Nok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
	}
}

func TestMatchICMPv4Event(t *testing.T) {
	ruleFilename := "icmpv4_rules.yml"
	pcapFilename := "icmpv4_values.pcap"
	rawPackets := false
	var rule Rule

	ruleset, err := LoadRuleFile(ruleFilename)
	if err != nil {
		t.Error(err)
		return
	}

	filteredEvents, _, err := ReadPacketsFromPcap(pcapFilename, layers.IPProtocolICMPv4, rawPackets)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filteredEvents) == 0 {
		t.Error("No ICMPv4 packets has been read from pcap")
		return
	}

	tests := []struct {
		Ok     []string
		Nok    []string
		Packet events.Event
	}{
		{
			Ok: []string{
				"ok_icmpv4_typecode",
				"ok_icmpv4_code",
				"ok_icmpv4_type",
				"ok_checksum",
			},
			Nok: []string{
				"nok_icmpv4_typecode",
				"nok_icmpv4_code",
				"nok_icmpv4_type",
				"nok_checksum",
			},
			Packet: filteredEvents[0],
		},
	}

	for _, suite := range tests {
		for _, rulename := range suite.Ok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); !ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
		for _, rulename := range suite.Nok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
	}
}

func TestMatchICMPv6Event(t *testing.T) {
	ruleFilename := "icmpv6_rules.yml"
	pcapFilename := "icmpv6_values.pcap"
	rawPackets := false
	var rule Rule

	ruleset, err := LoadRuleFile(ruleFilename)
	if err != nil {
		t.Error(err)
		return
	}

	filteredEvents, _, err := ReadPacketsFromPcap(pcapFilename, layers.IPProtocolICMPv6, rawPackets)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filteredEvents) == 0 {
		t.Error("No ICMPv6 packets has been read from pcap")
		return
	}

	tests := []struct {
		Ok     []string
		Nok    []string
		Packet events.Event
	}{
		{
			Ok: []string{
				"ok_icmpv6_typecode",
				"ok_icmpv6_code",
				"ok_icmpv6_type",
				"ok_checksum",
			},
			Nok: []string{
				"nok_icmpv6_typecode",
				"nok_icmpv6_code",
				"nok_icmpv6_type",
				"nok_checksum",
			},
			Packet: filteredEvents[0],
		},
	}

	for _, suite := range tests {
		for _, rulename := range suite.Ok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); !ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
		for _, rulename := range suite.Nok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
	}
}

func TestFilterIPs(t *testing.T) {
	ruleFilename := "ip_rules.yml"
	pcapFilename := "ip_values.pcap"
	rawPackets := false
	var rule Rule

	ruleset, err := LoadRuleFile(ruleFilename)
	if err != nil {
		t.Error(err)
		return
	}

	// use tcp packets in the capture to test IP filtering
	filteredEvents, _, err := ReadPacketsFromPcap(pcapFilename, layers.IPProtocolTCP, rawPackets)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filteredEvents) == 0 {
		t.Error("No TCP packets has been read from pcap")
		return
	}

	tests := []struct {
		Ok     []string
		Nok    []string
		Packet events.Event
	}{
		{
			Ok: []string{
				"ok_src_ips",
				"ok_src_ips_range",
			},
			Nok: []string{
				"nok_src_ips",
				"nok_src_ips_range",
			},
			Packet: filteredEvents[1],
		},
	}

	for _, suite := range tests {
		for _, rulename := range suite.Ok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); !ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
		for _, rulename := range suite.Nok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
	}
}

func TestFilterPorts(t *testing.T) {
	ruleFilename := "port_rules.yml"
	pcapFilename := "port_values.pcap"
	rawPackets := false
	var rule Rule

	ruleset, err := LoadRuleFile(ruleFilename)
	if err != nil {
		t.Error(err)
		return
	}

	// use tcp packets in the capture to test IP filtering
	filteredEvents, _, err := ReadPacketsFromPcap(pcapFilename, layers.IPProtocolTCP, rawPackets)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filteredEvents) == 0 {
		t.Error("No TCP packets has been read from pcap")
		return
	}

	tests := []struct {
		Ok     []string
		Nok    []string
		Packet events.Event
	}{
		{
			Ok: []string{
				"ok_dst_ports",
				"ok_dst_ports_range",
			},
			Nok: []string{
				"nok_dst_ports",
				"nok_dst_ports_range",
			},
			Packet: filteredEvents[1],
		},
	}

	for _, suite := range tests {
		for _, rulename := range suite.Ok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); !ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
		for _, rulename := range suite.Nok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
	}
}

func TestMatchTCPEvent(t *testing.T) {
	ruleFilename := "tcp_rules.yml"
	pcapFilename := "tcp_values.pcap"
	rawPackets := false
	var rule Rule

	ruleset, err := LoadRuleFile(ruleFilename)
	if err != nil {
		t.Error(err)
		return
	}

	filteredEvents, _, err := ReadPacketsFromPcap(pcapFilename, layers.IPProtocolTCP, rawPackets)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filteredEvents) == 0 {
		t.Error("No TCP packets has been read from pcap")
		return
	}

	tests := []struct {
		Ok     []string
		Nok    []string
		Packet events.Event
	}{
		{
			Ok: []string{
				"ok_ack",
				"ok_seq",
				"ok_window",
			},
			Nok: []string{
				"nok_ack",
				"nok_seq",
				"nok_window",
			},
			Packet: filteredEvents[1],
		},
		{
			Ok: []string{
				"ok_dsize",
				"ok_payload",
			},
			Nok: []string{
				"nok_dsize",
				"nok_payload",
			},
			Packet: filteredEvents[2],
		},
		{
			Ok: []string{
				"ok_flags",
				//"ok_fragbits",
			},
			Nok: []string{
				"nok_flags",
				//"nok_fragbits",
			},
			Packet: filteredEvents[4],
		},
	}

	for _, suite := range tests {
		for _, rulename := range suite.Ok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); !ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}

		for _, rulename := range suite.Nok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
	}
}

func TestMatchHTTPEvent(t *testing.T) {
	ruleFilename := "http_rules.yml"
	pcapFilename := "http_values.pcap"
	var httpEvents []*events.HTTPEvent
	var rule Rule

	// config.Cfg is initialized at the top of this file in the init function
	// Allow POST data to be parsed without loading the config.yml file
	config.Cfg.MaxPOSTDataSize = 8191

	ruleset, err := LoadRuleFile(ruleFilename)
	if err != nil {
		t.Error(err)
		return
	}

	packets, err := ReadRawTCPPacketsFromPcap(pcapFilename)
	if err != nil {
		t.Error(err)
		return
	}

	// Will fail on packets needing reassembly
	// Good enough for testing
	for _, packet := range packets {
		if len(packet.TransportLayer().LayerPayload()) > 0 {
			req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(packet.TransportLayer().LayerPayload())))
			if err != nil {
				t.Error(err)
				return
			}

			ev, _ := events.NewHTTPEvent(req, packet.NetworkLayer().NetworkFlow(), packet.TransportLayer().TransportFlow())
			if len(ev.Errors) > 0 {
				for _, err := range ev.Errors {
					t.Error(err)
				}
				return
			}
			httpEvents = append(httpEvents, ev)
		}
	}

	if httpEvents == nil {
		t.Error("No HTTP packets has been read from pcap")
		return
	}

	tests := []struct {
		Ok     []string
		Nok    []string
		Packet events.Event
	}{
		{
			Ok: []string{
				"ok_uri",
				"ok_body",
				"ok_is_tls",
				"ok_headers",
				"ok_proto",
				"ok_method",
			},
			Nok: []string{
				"nok_uri",
				"nok_body",
				"nok_is_tls",
				"nok_headers",
				"nok_proto",
				"nok_method",
			},
			Packet: httpEvents[0],
		},
	}

	for _, suite := range tests {
		for _, rulename := range suite.Ok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); !ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
		for _, rulename := range suite.Nok {
			rule = ruleset[rulename]
			if ok := rule.Match(suite.Packet); ok {
				t.Error(rulename, "FAILED")
				t.Fail()
			}
		}
	}
}
