package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ma111e/melody/internal/fileutils"

	"github.com/ma111e/melody/internal/config"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var (
	RootCmd = &cobra.Command{
		Use:   "meloctl",
		Short: "Melody helper",
		Long:  "Melody helper",
	}
	meloctlConf     *MeloctlConfig
	melodyConf      *config.Config
	meloctlConfFile string
	meloctlConfDir  string
)

func init() {
	// Trick to prevent the "unused variable" warning message for melodyConf, which is triggered because the relationship
	// between the cobra.OnInitialize function and the use of melodyConf via loadConfig is not being currently detected
	_ = melodyConf
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Failed to get home directory path")
	}

	meloctlConfDir = filepath.Join(homeDir, ".config", "meloctl")
	meloctlConfFile = filepath.Join(meloctlConfDir, "meloctl.yml")

	cobra.OnInitialize(loadConfig)
}

func loadConfig() {
	var err error
	meloctlConf = loadMeloctl()
	melodyConf, err = loadMelodyConfig(filepath.Join(meloctlConf.MelodyHomeDir, "config.yml"))
	if err != nil {
		log.Println(err)
	}
}

func loadMeloctl() *MeloctlConfig {
	var conf MeloctlConfig

	ok, err := fileutils.Exists(meloctlConfFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if !ok {
		if os.Args[1] != "init" {
			fmt.Println("⚠️ Meloctl config have not been initialized yet. Create a new configuration file using \"meloctl init\"")
		}

		return &conf
	}

	rawConf, err := ioutil.ReadFile(meloctlConfFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = yaml.Unmarshal(rawConf, &conf)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &conf
}

func loadMelodyConfig(configPath string) (*config.Config, error) {
	conf := config.NewConfig()
	cfgData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(cfgData, conf); err != nil {
		return nil, err
	}

	return conf, err
}
