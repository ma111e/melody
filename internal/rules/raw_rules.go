package rules

import (
	"fmt"
	"os"
	"strings"

	"github.com/ma111e/melody/internal/logging"

	"gopkg.in/yaml.v3"
)

// RawRules abstracts a group of raw rules in a rule file
type RawRules map[string]RawRule

// HTTPRule describes the raw "match" section of a rule targeting HTTP
type HTTPRule struct {
	URI     RawConditions `yaml:"http.uri"`
	Body    RawConditions `yaml:"http.body"`
	Headers RawConditions `yaml:"http.headers"`
	Verb    RawConditions `yaml:"http.method"`
	Proto   RawConditions `yaml:"http.proto"`
	TLS     *bool         `yaml:"http.tls"`
	Any     bool          `yaml:"any"`
}

// ParsedHTTPRule describes the parsed "match" section of a rule targeting HTTP
type ParsedHTTPRule struct {
	URI     *ConditionsList
	Body    *ConditionsList
	Headers *ConditionsList
	Verb    *ConditionsList
	Proto   *ConditionsList
	TLS     *bool
}

// TCPRule describes the raw "match" section of a rule targeting TCP
type TCPRule struct {
	IPOption RawConditions   `yaml:"tcp.ipoption"`
	Fragbits RawFragbitsList `yaml:"tcp.fragbits"`
	Dsize    *uint           `yaml:"tcp.dsize"`
	Flags    RawTCPFlagsList `yaml:"tcp.flags"`
	Seq      *uint32         `yaml:"tcp.seq"`
	Ack      *uint32         `yaml:"tcp.ack"`
	Payload  RawConditions   `yaml:"tcp.payload"`
	Window   *uint16         `yaml:"tcp.window"`
	Any      bool            `yaml:"any"`
}

// ParsedTCPRule describes the parsed "match" section of a rule targeting TCP
type ParsedTCPRule struct {
	IPOption *ConditionsList
	Fragbits []*uint8
	Flags    []*uint8
	Dsize    *uint
	Seq      *uint32
	Ack      *uint32
	Window   *uint16
	Payload  *ConditionsList
}

// ICMPv4Rule describes the raw "match" section of a rule targeting ICMPv4
type ICMPv4Rule struct {
	TypeCode *uint16       `yaml:"icmpv4.typecode"`
	Type     *uint8        `yaml:"icmpv4.type"`
	Code     *uint8        `yaml:"icmpv4.code"`
	Checksum *uint16       `yaml:"icmpv4.checksum"`
	Seq      *uint16       `yaml:"icmpv4.seq"`
	Payload  RawConditions `yaml:"icmpv4.payload"`
	Any      bool          `yaml:"any"`
}

// ParsedICMPv4Rule describes the parsed "match" section of a rule targeting ICMPv4
type ParsedICMPv4Rule struct {
	TypeCode *uint16
	Type     *uint8
	Code     *uint8
	Checksum *uint16
	Seq      *uint16
	Payload  *ConditionsList
}

// ICMPv6Rule describes the raw "match" section of a rule targeting ICMPv6
type ICMPv6Rule struct {
	TypeCode *uint16       `yaml:"icmpv6.typecode"`
	Type     *uint8        `yaml:"icmpv6.type"`
	Code     *uint8        `yaml:"icmpv6.code"`
	Checksum *uint16       `yaml:"icmpv6.checksum"`
	Payload  RawConditions `yaml:"icmpv6.payload"`
	Any      bool          `yaml:"any"`
}

// ParsedICMPv6Rule describes the parsed "match" section of a rule targeting ICMPv6
type ParsedICMPv6Rule struct {
	TypeCode *uint16
	Type     *uint8
	Code     *uint8
	Checksum *uint16
	Payload  *ConditionsList
}

// UDPRule describes the raw "match" section of a rule targeting UDP
type UDPRule struct {
	Length   *uint16       `yaml:"udp.length"`
	Dsize    *uint         `yaml:"udp.dsize"`
	Checksum *uint16       `yaml:"udp.checksum"`
	Payload  RawConditions `yaml:"udp.payload"`
	Any      bool          `yaml:"any"`
}

// ParsedUDPRule describes the parsed "match" section of a rule targeting UDP
type ParsedUDPRule struct {
	Length   *uint16
	Dsize    *uint
	Checksum *uint16
	Payload  *ConditionsList
}

// Filters groups the exposed rule filters
type Filters struct {
	Ports []string `yaml:"ports"`
	IPs   []string `yaml:"ips"`
}

// Metadata describes the exposed content of the "meta" field
type Metadata struct {
	ID          string   `yaml:"id"`
	Status      string   `yaml:"status"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Created     string   `yaml:"created"`
	Modified    string   `yaml:"modified"`
	References  []string `yaml:"references"`
}

// RawRule describes the format of a rule as written by the user
type RawRule struct {
	Whitelist Filters `yaml:"whitelist"`
	Blacklist Filters `yaml:"blacklist"`

	Match interface{} `yaml:"match"`

	Tags       map[string]string `yaml:"tags"`
	Layer      string            `yaml:"layer"`
	IPProtocol RawConditions     `yaml:"ip_protocol"`

	Metadata   Metadata          `yaml:"meta"`
	Additional map[string]string `yaml:"embed"`
}

var (
	validMatchKeysMap map[string]interface{} = LoadValidMatchKeysMap()
)

// Parse creates a Rule from a RawRule
func (rawRule RawRule) Parse() (Rule, error) {
	var err error
	rule := NewRule(rawRule)

	if rawRule.Match == nil {
		return rule, nil
	}

	rawMatch, err := yaml.Marshal(rawRule.Match)
	if err != nil {
		logging.Errors.Printf("failed to parse rule '%s' : invalid yaml definition (%s)", rawRule.Metadata.ID, err)
		// Fatal error
		os.Exit(1)
	}

	for key := range rawRule.Match.(map[string]interface{}) {
		// "any" is the only valid non-layer specific property in the "match" block :
		// It can be either "any" or <layer>.<property>
		if key != "any" && !strings.HasPrefix(key, rawRule.Layer+".") {
			return Rule{}, fmt.Errorf("property '%s' is not supported with layer '%s'", key, rawRule.Layer)
		}

		if _, ok := validMatchKeysMap[key]; !ok {
			return Rule{}, fmt.Errorf("unknown property '%s'", key)
		}
	}

	switch rawRule.Layer {
	case "http":
		var buf HTTPRule

		err = yaml.Unmarshal(rawMatch, &buf)
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedURI, err := buf.URI.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedBody, err := buf.Body.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedVerb, err := buf.Verb.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedHeaders, err := buf.Headers.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedProto, err := buf.Proto.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		rule.HTTP = ParsedHTTPRule{
			URI:     parsedURI,
			Body:    parsedBody,
			Verb:    parsedVerb,
			Headers: parsedHeaders,
			Proto:   parsedProto,
			TLS:     buf.TLS,
		}

		rule.MatchAll = !buf.Any

	case "tcp":
		var buf TCPRule

		err = yaml.Unmarshal(rawMatch, &buf)
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedIPOption, err := buf.IPOption.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedPayload, err := buf.Payload.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		rule.TCP = ParsedTCPRule{
			IPOption: parsedIPOption,
			Fragbits: buf.Fragbits.ParseList(),
			Flags:    buf.Flags.ParseList(),
			Window:   buf.Window,
			Dsize:    buf.Dsize,
			Seq:      buf.Seq,
			Ack:      buf.Ack,
			Payload:  parsedPayload,
		}

		rule.MatchAll = !buf.Any

	case "udp":
		var buf UDPRule

		err = yaml.Unmarshal(rawMatch, &buf)
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedPayload, err := buf.Payload.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		rule.UDP = ParsedUDPRule{
			Dsize:    buf.Dsize,
			Length:   buf.Length,
			Checksum: buf.Checksum,
			Payload:  parsedPayload,
		}

		rule.MatchAll = !buf.Any

	case "icmpv4":
		var buf ICMPv4Rule

		err = yaml.Unmarshal(rawMatch, &buf)
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedPayload, err := buf.Payload.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		rule.ICMPv4 = ParsedICMPv4Rule{
			TypeCode: buf.TypeCode,
			Type:     buf.Type,
			Code:     buf.Code,
			Checksum: buf.Checksum,
			Seq:      buf.Seq,
			Payload:  parsedPayload,
		}

		rule.MatchAll = !buf.Any

	case "icmpv6":
		var buf ICMPv6Rule

		err = yaml.Unmarshal(rawMatch, &buf)
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		parsedPayload, err := buf.Payload.ParseList()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		}

		rule.ICMPv6 = ParsedICMPv6Rule{
			TypeCode: buf.TypeCode,
			Type:     buf.Type,
			Code:     buf.Code,
			Checksum: buf.Checksum,
			Payload:  parsedPayload,
		}

		rule.MatchAll = !buf.Any
	}

	return rule, nil
}
