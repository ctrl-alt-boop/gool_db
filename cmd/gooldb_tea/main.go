package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	tui "github.com/ctrl-alt-boop/gooldb/tea"
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

	tui := tui.NewModel(gool)
	p := tea.NewProgram(tui, tea.WithAltScreen())
	tui.SetProgramSend(p.Send)
	if _, err := p.Run(); err != nil {
		fmt.Printf("GoolTea error: %v", err)
		os.Exit(1)
	}
}
