package main

import (
	"fmt"
	"os"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	"github.com/ctrl-alt-boop/gooldb/tui"
	"github.com/jesseduffield/gocui"
)

func main() {
	var ip string = "localhost"
	if len(os.Args) > 1 {
		ip = os.Args[1]
	}

	logger := logging.NewLogger("app.log")
	defer logger.Close()

	notifier := tui.NewNotifier()

	logger.Info("GoolDb Create")
	gool := gooldb.Create(logger, notifier, ip)

	logger.Info("Tui Create")
	tui := tui.Create(notifier, gool)

	logger.Info("Tui Run")
	err := panicRecovery(func() error {
		return tui.Run()
	}, logger)

	if err != nil && err != gocui.ErrQuit {
		logger.Panic(err)
	}
}

func panicRecovery(fn func() error, logger *logging.Logger) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = fmt.Errorf("panicRecovery: %w", e)
			} else {
				err = fmt.Errorf("panicRecovery: %v", r)
			}
			logger.Panic(err)
		}
	}()
	return fn()
}
