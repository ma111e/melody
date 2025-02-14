package rules

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/ma111e/melody/internal/events"
	"os"
)

var (
	assetsBasePath string
)

func init() {
	assetsBasePath = "test_resources"
}

// ReadRawTCPPacketsFromPcap is an helper that reads raw TCP packets from the specified pcap file
func ReadRawTCPPacketsFromPcap(pcapfile string) ([]gopacket.Packet, error) {
	var packets []gopacket.Packet
	_, rawPackets, err := ReadPacketsFromPcap(pcapfile, layers.IPProtocolTCP, true)
	if err != nil {
		return nil, err
	}

	for _, val := range rawPackets {
		packets = append(packets, val.(gopacket.Packet))
	}
	return packets, nil
}

// ReadPacketsFromPcap is an helper that reads packets from the specified pcap file and returns them as an array of Event
func ReadPacketsFromPcap(pcapfile string, filter layers.IPProtocol, raw bool) ([]events.Event, []gopacket.Packet, error) {
	var Events []events.Event
	var rawPackets []gopacket.Packet
	var ret []events.Event
	var rawRet []gopacket.Packet
	pcapfilePath := MakeAssetFullPath(pcapfile)

	f, err := os.Open(pcapfilePath)
	if err != nil {
		return []events.Event{}, []gopacket.Packet{}, err
	}
	handle, err := pcap.OpenOfflineFile(f)
	if err != nil {
		return []events.Event{}, []gopacket.Packet{}, err
	}

	src := gopacket.NewPacketSource(handle, handle.LinkType())
	in := src.Packets()

loop:
	for {
		packet := <-in
		if packet == nil {
			break loop
		}

		if _, ok := packet.NetworkLayer().(*layers.IPv4); ok {
			if packet.NetworkLayer().(*layers.IPv4).Protocol == filter {
				if raw {
					rawPackets = append(rawPackets, packet)
				} else {
					switch filter {
					case layers.IPProtocolICMPv4:
						ev, err := events.NewICMPv4Event(packet)
						if err != nil {
							return []events.Event{}, []gopacket.Packet{}, err
						}

						Events = append(Events, ev)

					case layers.IPProtocolUDP:
						ev, err := events.NewUDPEvent(packet, 4)
						if err != nil {
							return []events.Event{}, []gopacket.Packet{}, err
						}

						Events = append(Events, ev)

					case layers.IPProtocolTCP:
						ev, err := events.NewTCPEvent(packet, 4)
						if err != nil {
							return []events.Event{}, []gopacket.Packet{}, err
						}

						Events = append(Events, ev)

					default:
						continue loop
					}
				}
			}
		} else if _, ok := packet.NetworkLayer().(*layers.IPv6); ok {
			if packet.NetworkLayer().(*layers.IPv6).NextHeader == filter {
				switch filter {
				case layers.IPProtocolICMPv6:
					ev, err := events.NewICMPv6Event(packet)
					if err != nil {
						return []events.Event{}, []gopacket.Packet{}, err
					}

					Events = append(Events, ev)

				default:
					continue loop
				}
			}
		}
	}

	// I'm so lazy
	if raw {
		rawRet = make([]gopacket.Packet, len(rawPackets))
		copy(rawRet, rawPackets)
	}

	ret = make([]events.Event, len(Events))
	copy(ret, Events)

	return ret, rawRet, nil
}

// LoadRuleFile is an helper that parses the rule file at the given path and returns a rule set
func LoadRuleFile(rulefile string) (map[string]Rule, error) {
	ruleset := make(map[string]Rule)
	rulefilePath := MakeAssetFullPath(rulefile)
	rawRules, err := ParseYAMLRulesFile(rulefilePath)
	if err != nil {
		return map[string]Rule{}, err
	}
	for name, rawRule := range rawRules {
		ruleset[name], err = rawRule.Parse()
		if err != nil {
			return ruleset, err
		}
	}

	return ruleset, nil
}

// MakeAssetFullPath is an helper that returns the path to the tests resources as defined by the assetsBasePath variable
func MakeAssetFullPath(path string) string {
	return fmt.Sprintf("%s/%s", assetsBasePath, path)
}
