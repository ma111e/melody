package main

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/ma111e/melody/internal/config"
	"github.com/ma111e/melody/internal/tagparser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	setCmd = &cobra.Command{
		Args:  cobra.ExactArgs(2),
		Use:   "set",
		Short: "Set a Meloctl config value by name",
		Long:  `This subcommand is used to set a Meloctl config value by name`,
		Run:   setConfigKey,
	}
)

func init() {
	RootCmd.AddCommand(setCmd)
}

func setConfigKey(_ *cobra.Command, args []string) {
	var found bool
	key := args[0]
	val := args[1]

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

	if err := setValueByYAMLTagName(meloctlConf, key, val); err != nil {
		fmt.Printf("❌ Error while setting value : %s\n", err)
		return
	}

	fmt.Printf("%s => %s \n", key, val)

	out, err := yaml.Marshal(meloctlConf)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(meloctlConfFile, out, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("✅ [%s] Configuration file updated\n", meloctlConfFile)
}

func setValueByYAMLTagName(what interface{}, key string, val string) error {
	for i := 0; i < reflect.TypeOf(what).Elem().NumField(); i++ {
		tag := reflect.TypeOf(what).Elem().Field(i).Tag

		_, exists := tag.Lookup("yaml")
		if !exists {
			continue
		}

		tagValue, err := tagparser.ParseYamlTagValue(tag)
		if err != nil {
			return err
		}

		if key == tagValue {
			reflect.ValueOf(what).Elem().FieldByName(reflect.TypeOf(what).Elem().Field(i).Name).SetString(val)
			break
		}
	}

	return nil
}
