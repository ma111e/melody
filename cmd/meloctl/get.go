package main

import (
	"fmt"
	"github.com/ma111e/melody/internal/config"
	"github.com/ma111e/melody/internal/tagparser"
	"reflect"

	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "get",
		Short: "Get a Meloctl config value by name",
		Long:  `This subcommand is used to get a Meloctl config value by name`,
		Run:   getConfigKey,
	}
)

func init() {
	RootCmd.AddCommand(getCmd)
}

func getConfigKey(_ *cobra.Command, args []string) {
	var found bool
	key := args[0]

	validKeys, _ := config.LoadYAMLTagsOf(MeloctlConfig{})
	for _, valid := range validKeys {
		if key == valid {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("❌ [%s] Unknown key\n", key)
		return
	}

	val, err := getValueByYAMLTagName(meloctlConf, key)
	if err != nil {
		fmt.Printf("❌ Error while fetching value : %s\n", err)
		return
	}

	fmt.Printf("%s => %s \n", key, val)
}

func getValueByYAMLTagName(what interface{}, key string) (string, error) {
	var val string

	for i := 0; i < reflect.TypeOf(what).Elem().NumField(); i++ {
		tag := reflect.TypeOf(what).Elem().Field(i).Tag

		_, exists := tag.Lookup("yaml")
		if !exists {
			continue
		}

		tagValue, err := tagparser.ParseYamlTagValue(tag)
		if err != nil {
			return val, err
		}

		if key == tagValue {
			val = reflect.ValueOf(what).Elem().FieldByName(reflect.TypeOf(what).Elem().Field(i).Name).String()
			break
		}
	}

	return val, nil
}
