package logdata

import (
	"fmt"

	"github.com/ma111e/melody/internal/loggable"
)

// IPLogData is the interface used by packet structs supporting an IP layer
type IPLogData interface{}

// BaseLogData is used as the base packet log and contains common data, such as the timestamp
type BaseLogData struct {
	Timestamp  string              `json:"timestamp"`
	Session    string              `json:"session"`
	Type       string              `json:"type"`
	SourceIP   string              `json:"src_ip"`
	DestPort   uint16              `json:"dst_port"`
	Tags       map[string][]string `json:"matches"`
	InlineTags []string            `json:"inline_matches"`
	Additional map[string]string   `json:"embedded"`
}

// Init takes the common BaseEvent attributes to setup the BaseLogData struct
func (l *BaseLogData) Init(ev loggable.Loggable) {
	l.Type = ev.GetKind()
	l.SourceIP = ev.GetSourceIP()
	l.DestPort = ev.GetDestPort()
	l.Session = ev.GetSession()
	l.InlineTags = []string{}

	if len(ev.GetTags()) == 0 {
		l.Tags = make(map[string][]string)
	} else {
		l.Tags = ev.GetTags()
		l.InlineTags = makeInlineArray(ev.GetTags())
	}
}

// makeInlineArray converts a Tags map to an array of its values with the keys and values merged with a '.'
func makeInlineArray(tags map[string][]string) []string {
	var inlineTags []string

	for key, values := range tags {
		for _, val := range values {
			inlineTags = append(inlineTags, fmt.Sprintf("%s.%s", key, val))
		}
	}

	return inlineTags
}
