package rules

import (
	"reflect"

	"github.com/ma111e/melody/internal/tagparser"
)

// LoadValidMatchKeysMap returns a map of the json keys for each of the protos structs
func LoadValidMatchKeysMap() map[string]interface{} {
	loadFns := []func() ([]string, error){
		loadHTTPYamlTags,
		loadTCPYamlTags,
		loadUDPYamlTags,
		loadICMPv4YamlTags,
		loadICMPv6YamlTags,
	}

	matchKeysMap := make(map[string]interface{})
	for _, loadFn := range loadFns {
		tags, err := loadFn()
		if err != nil {
			panic(err)
		}

		for _, tag := range tags {
			matchKeysMap[tag] = new(interface{})
		}
	}

	return matchKeysMap
}

// Below : all the same functions with a different struct
func loadHTTPYamlTags() ([]string, error) {
	var tags []string
	for i := 0; i < reflect.TypeOf(HTTPRule{}).NumField(); i++ {
		ruleTag := reflect.TypeOf(HTTPRule{}).Field(i).Tag
		tagValue, err := tagparser.ParseYamlTagValue(ruleTag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tagValue)
	}

	return tags, nil
}

func loadTCPYamlTags() ([]string, error) {
	var tags []string
	for i := 0; i < reflect.TypeOf(TCPRule{}).NumField(); i++ {
		ruleTag := reflect.TypeOf(TCPRule{}).Field(i).Tag
		tagValue, err := tagparser.ParseYamlTagValue(ruleTag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tagValue)
	}

	return tags, nil
}

func loadUDPYamlTags() ([]string, error) {
	var tags []string
	for i := 0; i < reflect.TypeOf(UDPRule{}).NumField(); i++ {
		ruleTag := reflect.TypeOf(UDPRule{}).Field(i).Tag
		tagValue, err := tagparser.ParseYamlTagValue(ruleTag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tagValue)
	}

	return tags, nil
}

func loadICMPv4YamlTags() ([]string, error) {
	var tags []string
	for i := 0; i < reflect.TypeOf(ICMPv4Rule{}).NumField(); i++ {
		ruleTag := reflect.TypeOf(ICMPv4Rule{}).Field(i).Tag
		tagValue, err := tagparser.ParseYamlTagValue(ruleTag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tagValue)
	}

	return tags, nil
}

func loadICMPv6YamlTags() ([]string, error) {
	var tags []string
	for i := 0; i < reflect.TypeOf(ICMPv6Rule{}).NumField(); i++ {
		ruleTag := reflect.TypeOf(ICMPv6Rule{}).Field(i).Tag
		tagValue, err := tagparser.ParseYamlTagValue(ruleTag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tagValue)
	}

	return tags, nil
}
