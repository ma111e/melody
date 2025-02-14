package config

import (
	"reflect"

	"github.com/ma111e/melody/internal/tagparser"
)

// LoadYAMLTagsOf loads the yaml tags of a struct
func LoadYAMLTagsOf(what interface{}) ([]string, error) {
	var tags []string

	for i := 0; i < reflect.TypeOf(what).NumField(); i++ {
		ruleTag := reflect.TypeOf(what).Field(i).Tag

		_, exists := ruleTag.Lookup("yaml")
		if !exists {
			continue
		}

		tagValue, err := tagparser.ParseYamlTagValue(ruleTag)
		if err != nil {
			return tags, err
		}

		tags = append(tags, tagValue)
	}

	return tags, nil
}

// LoadValidConfigKeysMap returns a map of the json keys present in the Config struct
func LoadValidConfigKeysMap() map[string]interface{} {
	configKeysMap := make(map[string]interface{})
	tags, err := LoadYAMLTagsOf(Config{})
	if err != nil {
		panic(err)
	}

	for _, tag := range tags {
		configKeysMap[tag] = new(interface{})
	}

	return configKeysMap
}
