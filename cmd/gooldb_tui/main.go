package main

import (
	"log"
	"os"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/tui"
	"github.com/jroimartin/gocui"
)

func main() {
	var ip string = "localhost"
	if len(os.Args) > 1 {
		ip = os.Args[1]
	}
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatalf("failed to create log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gool := gooldb.Create(ip)

	tui := tui.Create(gool)

	if err := tui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
