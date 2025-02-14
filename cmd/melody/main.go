package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ma111e/melody/internal/engine"
	"github.com/ma111e/melody/internal/rules"
	"github.com/ma111e/melody/internal/sensor"

	"github.com/ma111e/melody/internal/logging"

	"github.com/google/shlex"
	"github.com/ma111e/melody/internal/config"
	"github.com/pborman/getopt/v2"
)

var (
	quitErrChan       = make(chan error)
	shutdownChan      = make(chan bool)
	loggerStoppedChan = make(chan bool)
	engineStoppedChan = make(chan bool)
	sensorStoppedChan = make(chan bool)
	quitSigChan       = make(chan os.Signal, 1)
)

func init() {
	signal.Notify(quitSigChan, syscall.SIGINT, syscall.SIGTERM)

	config.Cli.PcapFilePath = getopt.StringLong("pcap", 'r', "", "Replay a pcap file into the sensor")
	config.Cli.BPF = getopt.StringLong("filter", 'F', "", "Override the filter.bpf file with the specified filter")
	config.Cli.HomeDirPath = getopt.StringLong("home-dir", 'H', "", "Set the home directory")
	config.Cli.ConfigDirPath = getopt.StringLong("config-dir", 'C', "", "Set the config directory")
	config.Cli.ConfigFilePath = getopt.StringLong("config", 'c', "", "Path to the config file to load")
	config.Cli.BPFFilePath = getopt.StringLong("bpf", 'f', "", "Path to the BPF file")
	config.Cli.Interface = getopt.StringLong("interface", 'i', "", "Listen on the specified interface")
	config.Cli.Stdout = getopt.BoolLong("stdout", 's', "Output logged data to stdout instead")
	config.Cli.Dump = getopt.BoolLong("dump", 'd', "Output raw packet details instead of JSON")
	getopt.FlagLong(&config.Cli.FreeConfig, "option", 'o', "Override configuration keys")
	getopt.Parse()

	if line := os.Getenv("MELODY_CLI"); line != "" {
		chunks, err := shlex.Split(line)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to chunk MELODY_CLI env variable")
			os.Exit(1)
		}

		chunks = append([]string{os.Args[0]}, chunks...)
		getopt.CommandLine.Parse(chunks)
	}

	config.Cfg = config.NewConfig()
	err := config.Cfg.Load()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = logging.InitLoggers()
	if err != nil {
		logging.Std.Println(err)
		os.Exit(1)
	}
	loaded := rules.LoadRulesDir(filepath.Join(config.Cfg.HomeDirPath, config.Cfg.RulesDir))

	logging.Std.Printf("Loaded %d rules\n", loaded)
	logging.Std.Printf("Listing on interface %s\n", config.Cfg.Interface)
}

func main() {
	logging.Std.Println("Starting loggers")
	logging.Start(quitErrChan, shutdownChan, loggerStoppedChan)

	logging.Std.Println("Starting engine")
	engine.Start(quitErrChan, shutdownChan, engineStoppedChan)

	logging.Std.Println("Starting sensor")
	sensor.Start(quitErrChan, shutdownChan, sensorStoppedChan)

	logging.Std.Println("All systems started")

	select {
	case err := <-quitErrChan:
		logging.Errors.Println(err)
		close(shutdownChan)
		break
	case <-quitSigChan:
		close(shutdownChan)
		break
	case <-shutdownChan:
		logging.Std.Println("Shutting down")
		break
	}

	<-sensorStoppedChan
	logging.Std.Println("Sensor stopped")

	<-engineStoppedChan
	logging.Std.Println("Engine stopped")

	<-loggerStoppedChan
	logging.Std.Println("Logger stopped")

	logging.Std.Println("Reached shutdown")
}
