package main

import (
	"log"
	"os"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/tui/terminal"
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

	log.Println("Hello")
	tui := terminal.Create()
	log.Println("Bello")

	go func() {
		db, err := gooldb.Connect(ip)
		log.Println("Connect() success")
		if err != nil {
			log.Panicln(err)
		}
		log.Println("SetDatabaseContext()")
		tui.SetDatabaseContext(db)
		log.Println("SetDatabaseContext() done")
	}()
	log.Println("Dello")
	if err := tui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
