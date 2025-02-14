package rules

import (
	"github.com/ma111e/melody/internal/filters"
	"github.com/ma111e/melody/internal/logging"
)

// Rules abstracts an array of Rule
type Rules []Rule

// Rule describes a parsed Rule object, used to match against a byte array
type Rule struct {
	Name string
	ID   string
	Tags map[string]string

	Layer string

	IPProtocol *ConditionsList

	HTTP   ParsedHTTPRule
	TCP    ParsedTCPRule
	UDP    ParsedUDPRule
	ICMPv4 ParsedICMPv4Rule
	ICMPv6 ParsedICMPv6Rule

	IPs        filters.IPRules
	Ports      filters.PortRules
	Metadata   Metadata
	Additional map[string]string

	MatchAll bool
}

// NewRule creates a Rule from a RawRule
func NewRule(rawRule RawRule) Rule {
	var portsList = filters.PortRules{
		WhitelistedPorts: filters.PortRanges{},
		BlacklistedPorts: filters.PortRanges{},
	}

	var ipsList = filters.IPRules{
		WhitelistedIPs: filters.IPRanges{},
		BlacklistedIPs: filters.IPRanges{},
	}

	ipsList.ParseRules(rawRule.Whitelist.IPs, rawRule.Blacklist.IPs)
	portsList.ParseRules(rawRule.Whitelist.Ports, rawRule.Blacklist.Ports)

	parsedIPProtocol, err := rawRule.IPProtocol.ParseList()
	if err != nil {
		logging.Errors.Printf("failed to parse rule '%s' : %s", rawRule.Metadata.ID, err)
		return Rule{}
	}

	rule := Rule{
		Tags:       rawRule.Tags,
		IPProtocol: parsedIPProtocol,
		ID:         rawRule.Metadata.ID,
		Layer:      rawRule.Layer,
		Ports:      portsList,
		IPs:        ipsList,
		Metadata:   rawRule.Metadata,
		Additional: rawRule.Additional,
	}

	return rule
}

// Filter is a helper filtering out one or multiple Rule according to a function returning true or false
// Similar to array.filter() in python
func (rules Rules) Filter(fn func(rule Rule) bool) Rules {
	res := Rules{}

	for _, rule := range rules {
		if fn(rule) {
			res = append(res, rule)
		}
	}

	return res
}

// Parse parses raw rules to create a set of rules as Rules
//func (rawRules RawRules) Parse() Rules {
//	rules := Rules{}
//	for rname, rule := range rawRules {
//		parsedRule := rule.Parse()
//		parsedRule.Name = rname
//		rules = append(rules, parsedRule)
//	}
//
//	return rules
//}
