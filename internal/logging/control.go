package logging

import (
	"github.com/ma111e/melody/internal/config"
	"github.com/ma111e/melody/internal/events"
)

var (
	// LogChan is the channel used to receive events to be logged
	LogChan = make(chan events.Event)
)

// Start starts the logging pipeline
func Start(quitErrChan chan error, shutdownChan chan bool, loggerStoppedChan chan bool) {
	go receiveEventsForLogging(quitErrChan, shutdownChan, loggerStoppedChan)
}

func isIPv4(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			return true
		}
	}
	return false
}

func receiveEventsForLogging(quitErrChan chan error, shutdownChan chan bool, loggerStoppedChan chan bool) {
	defer func() {
		close(loggerStoppedChan)
	}()

	for {
		select {

		case ev := <-LogChan:
			switch ev.GetKind() {
			case config.HTTPKind:
				if isIPv4(ev.GetSourceIP()) {
					if _, ok := config.Cfg.DiscardProto4[config.HTTPKind]; ok {
						continue
					}
				} else {
					if _, ok := config.Cfg.DiscardProto6[config.HTTPKind]; ok {
						continue
					}
				}
			case config.HTTPSKind:
				if isIPv4(ev.GetSourceIP()) {
					if _, ok := config.Cfg.DiscardProto4[config.HTTPSKind]; ok {
						continue
					}
				} else {
					if _, ok := config.Cfg.DiscardProto6[config.HTTPSKind]; ok {
						continue
					}
				}
			}
			logdata, err := ev.ToLog().String()
			if err != nil {
				Warnings.Println("Failed to serialize JSON payload while writing to log file")
				continue
			}

			// Log to sensor file
			Sensor.Println(logdata)

		case <-quitErrChan:
			return

		case <-shutdownChan:
			return
		}
	}
}
