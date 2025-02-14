package rules

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// RawConditions describes the format of a condition field in a rule file
type RawConditions struct {
	Groups map[string][]string `yaml:"-,inline"`
	Any    bool                `yaml:"any"`
	Depth  uint                `yaml:"depth"`
	Offset uint                `yaml:"offset"`
}

// ConditionsList describes the format of a list of RawConditions
type ConditionsList struct {
	Conditions []Conditions
	MatchAll   bool
	//MatchAny   bool
}

// Conditions describes a parsed RawConditions
type Conditions struct {
	Values  []ConditionValue
	Options Options
}

// Options describes the available matching options
type Options struct {
	Depth      uint
	Offset     uint
	Nocase     bool
	Is         bool
	All        bool
	Contains   bool
	Startswith bool
	Endswith   bool
	Regex      bool
}

// ConditionValue abstracts the parsed value of a condition to use in a match attempt
type ConditionValue struct {
	CompiledRegex *regexp.Regexp
	ByteValue     []byte
}

// Match matches a byte array against a ConditionsList
func (clst ConditionsList) Match(received []byte) bool {
	if clst.Conditions != nil {
		if !clst.MatchAll {
			var condOK = false
			for _, condGroup := range clst.Conditions {
				// If any condition group is valid, continue

				if condGroup.Match(received) {
					condOK = true
					break
				}
			}

			if !condOK {
				return false
			}
		} else { // condList.MatchAll
			for _, condGroup := range clst.Conditions {
				// If any condition group is invalid, rule is false
				// Continue if the test for all the values are successful

				if !condGroup.Match(received) {
					return false
				}
			}
		}
	}

	return true
}

// Match matches a byte array against a set of conditions
func (cds Conditions) Match(received []byte) bool {
	var contentMatch bool
	var matchCounter int
	var valuesLen = len(cds.Values)

	for _, condVal := range cds.Values {
		contentMatch = cds.MatchBytesWithOptions(received, condVal)

		if cds.Options.All && !contentMatch {
			return false
			//	Continue unless all the tests passed
		} else if cds.Options.All && matchCounter < valuesLen {
			matchCounter++
			continue
		} else if contentMatch {
			// If at least one match in the list, then continue
			break
		}
	}

	return contentMatch
}

// MatchBytesWithOptions matches a byte array against a set of conditions, according to the specified ConditionValue
// This function only cares about the matching modifier ("contains", "startswith", etc, not "all")
// The condition's options are being taken care of in the Conditions.Match function
func (cds Conditions) MatchBytesWithOptions(received []byte, condVal ConditionValue) bool {
	var match bool
	var condValContent = condVal.ByteValue

	if cds.Options.Nocase {
		received = bytes.ToLower(received)
		condValContent = bytes.ToLower(condValContent)
	}

	if cds.Options.Offset > 0 && cds.Options.Offset < uint(len(received)) {
		received = received[cds.Options.Offset:]
	}

	if cds.Options.Depth > 0 && cds.Options.Depth < uint(len(received)) {
		received = received[:cds.Options.Depth]
	}

	if cds.Options.Is {
		match = bytes.Equal(received, condValContent)
	} else if cds.Options.Regex {
		match = condVal.CompiledRegex.Match(received)
	} else if cds.Options.Contains {
		match = bytes.Contains(received, condValContent)
	} else if cds.Options.Startswith {
		match = bytes.HasPrefix(received, condValContent)
	} else if cds.Options.Endswith {
		match = bytes.HasSuffix(received, condValContent)
	}
	return match
}

// ParseList parses a RawConditions set to create a ConditionsList
func (rclst RawConditions) ParseList() (*ConditionsList, error) {
	if len(rclst.Groups) == 0 {
		return nil, nil
	}

	condsList := ConditionsList{
		MatchAll: !rclst.Any,
	}
	var bufCond Conditions

	for options, val := range rclst.Groups {
		bufCond = Conditions{}
		err := bufCond.ParseOptions(options)
		if err != nil {
			return nil, err
		}
		bufCond.ParseValues(val)
		bufCond.Options.Offset = rclst.Offset
		bufCond.Options.Depth = rclst.Depth
		//bufCond.Options.All = rclst.Any == false

		condsList.Conditions = append(condsList.Conditions, bufCond)
	}

	//condsList.ParseMatchType(rclst.MatchType, ruleID)

	return &condsList, nil
}

// ParseOptions parses a condition's name to extract the options separated by a |
func (cds *Conditions) ParseOptions(opt string) error {
	chunks := strings.Split(opt, "|")
	modeQty := 0
	var newOption Options

	// Default to all = true
	newOption.All = true

	if opt == "" {
		return fmt.Errorf("options httpparser failed for condition '%s' : matching mode cannot be empty", opt)
	}

	for _, chunk := range chunks {
		switch chunk {
		case "any":
			newOption.All = false
		case "nocase":
			newOption.Nocase = true
		case "regex":
			newOption.Regex = true
		case "is":
			modeQty++
			newOption.Is = true
		case "contains":
			modeQty++
			newOption.Contains = true
		case "startswith":
			modeQty++
			newOption.Startswith = true
		case "endswith":
			modeQty++
			newOption.Endswith = true
		default:
			return fmt.Errorf("options httpparser failed for condition '%s' : unknown option \"%s\"", opt, chunk)
		}
	}

	if modeQty > 1 {
		return fmt.Errorf("options httpparser failed for condition '%s' : there can only be one of <nocase|regex|is|contains|startswith|endswith>", opt)
	}

	//newOption.All = any == false
	cds.Options = newOption

	return nil
}

// ParseValues loads a Conditions set from a list of condition strings
func (cds *Conditions) ParseValues(list []string) {
	var err error
	var condValBuf = ConditionValue{}

	for _, val := range list {
		buffer := []byte(val)

		condValBuf = ConditionValue{}
		if cds.Options.Regex {
			condValBuf.CompiledRegex, err = regexp.Compile(val)
			if err != nil {
				log.Println("Failed to compile regex", val, ":", err)
				os.Exit(1)
			}
		}
		parsedBuffer, err := ParseHybridPattern(buffer)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		condValBuf.ByteValue = parsedBuffer

		cds.Values = append(cds.Values, condValBuf)
	}
}

// ParseHybridPattern parses a byte array composed of hybrid hex and ascii characters and returns its
// equivalent as a byte array
func ParseHybridPattern(buffer []byte) ([]byte, error) {
	var isHex bool
	var parsedBuffer []byte
	var byteBuffer []byte

	for _, c := range buffer {
		if c == byte('|') {
			// If we already met a '|', then this one is the end delimiter
			// -> Dump the recorded byte buffer and clean it
			if isHex {
				isHex = false

				data, err := hex.DecodeString(string(byteBuffer))
				if err != nil {
					return nil, fmt.Errorf("failed to parse hybrid pattern : [%s] in %s", err, buffer)
				}

				parsedBuffer = append(parsedBuffer, data...)
				byteBuffer = []byte{}
				continue
			} else {
				// Else start recording the bytes until a delimiter is found
				isHex = true
				continue
			}
		} else if isHex {
			// If we already have met a byte delimiter and this char is not a delimiter as well, record as byte
			// Skip spaces in user defined hex string
			if c == ' ' {
				continue
			}
			byteBuffer = append(byteBuffer, c)
		} else {
			parsedBuffer = append(parsedBuffer, c)
		}
	}

	if isHex {
		return nil, fmt.Errorf("failed to parse hybrid pattern : uneven number of hex delimiter (\"|\") in %s", buffer)
	}

	return parsedBuffer, nil
}
