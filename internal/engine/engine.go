package engine

import (
	"github.com/ma111e/melody/internal/config"
	"github.com/ma111e/melody/internal/events"
	"github.com/ma111e/melody/internal/logging"
	"github.com/ma111e/melody/internal/router"
	"github.com/ma111e/melody/internal/rules"
)

var (
	// EventChan is the channel used to receive event to qualify
	EventChan = make(chan events.Event)
)

// Start starts the matching and tagging engine
func Start(quitErrChan chan error, shutdownChan chan bool, engineStoppedChan chan bool) {
	go startEventQualifier(quitErrChan, shutdownChan, engineStoppedChan)

	if config.Cfg.ServerHTTPEnable {
		logging.Std.Println("Starting HTTP server")
		go router.StartHTTP(quitErrChan)
	}

	if config.Cfg.ServerHTTPSEnable {
		logging.Std.Println("Starting HTTPS server")
		go router.StartHTTPS(quitErrChan, EventChan)
	}
}

func startEventQualifier(quitErrChan chan error, shutdownChan chan bool, engineStoppedChan chan bool) {
	defer func() {
		close(engineStoppedChan)
	}()

	for {
		select {
		case <-shutdownChan:
			return

		case <-quitErrChan:
			return

		case ev := <-EventChan:
			var matches []rules.Rule

			for _, ruleset := range rules.GlobalRules[ev.GetKind()] {
				for _, rule := range ruleset {
					if rule.Match(ev) {
						matches = append(matches, rule)
					}
				}
			}

			if len(matches) > 0 {
				for _, match := range matches {
					ev.AddTags(match.Tags)
					ev.AddAdditional(match.Additional)
				}
			}

			logging.LogChan <- ev
		}
	}
}
