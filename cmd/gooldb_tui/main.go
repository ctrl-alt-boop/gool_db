package main

import (
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
	if err := tui.MainLoop(); err != nil && err != gocui.ErrQuit {
		logger.Panic(err)
	}
}
