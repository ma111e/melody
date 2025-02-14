package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/ma111e/melody/internal/config"

	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Interact with a Melody config file",
		Long:  `This subcommand is used to interact with a Melody config file`,
	}
	checkConfigCmd = &cobra.Command{
		Use:   "check",
		Args:  cobra.ExactArgs(1),
		Short: "Validate a Melody config file",
		Long:  `This subcommand is used to validate a Melody config file`,
		Run:   checkConfig,
	}

	validConfigKeysMap map[string]interface{} = config.LoadValidConfigKeysMap()
)

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(checkConfigCmd)
}

func checkConfig(_ *cobra.Command, args []string) {
	var err error
	configPath := args[0]

	cfg := make(map[string]interface{})

	cfgData, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("❌ [%s]: %s\n", configPath, err)
		return
	}

	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		fmt.Printf("❌ [%s]: %s\n", configPath, err)
		return
	}

	for key := range cfg {
		if _, ok := validConfigKeysMap[key]; !ok {
			fmt.Printf("❌ [%s]: unknown property '%s'\n", configPath, key)
			return
		}
	}

	fmt.Printf("✅ [%s]: OK\n", configPath)
}
