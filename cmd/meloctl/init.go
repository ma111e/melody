package main

import (
	"fmt"
	"github.com/ma111e/melody/internal/fileutils"
	"github.com/ma111e/melody/internal/meloctl/prompt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

// MeloctlConfig represents the meloctl config file
type MeloctlConfig struct {
	MelodyHomeDir string `yaml:"melody.home" pretty:"Melody home directory"`
	//MelodyRulesSource string `yaml:"melody.rules.sources" pretty:"Melody rules sources"`
}

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Create Meloctl config",
		Long:  `This subcommand is used to setup Meloctl`,
		Run:   initConfig,
	}

	rawMeloconf = map[string]string{}
)

func init() {
	RootCmd.AddCommand(initCmd)
}

func initConfig(_ *cobra.Command, args []string) {
	_ = os.MkdirAll(meloctlConfDir, 0755)

	if ok, _ := fileutils.Exists(meloctlConfFile); ok {
		fmt.Printf("✅ [%s] Meloctl is already installed\n", meloctlConfFile)
		return
	}

	err := prompt.AskAll(MeloctlConfig{}, &rawMeloconf)
	if err != nil {
		fmt.Println(err)
		return
	}

	conf := NewMeloctlConfigFromRaw(rawMeloconf)
	out, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(meloctlConfFile, out, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("✅ [%s] Meloctl has been initialized\n", meloctlConfFile)
}

// NewMeloctlConfigFromRaw creates a default MeloctlConfig struct
func NewMeloctlConfigFromRaw(raw map[string]string) *MeloctlConfig {
	return &MeloctlConfig{
		MelodyHomeDir: raw["melody.home"],
		//MelodyRulesSource: raw["melody.rules.sources"],
	}
}
