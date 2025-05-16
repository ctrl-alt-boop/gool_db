package event

import "github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"

type GoolDbEventMsg struct {
	Type gooldb.EventType
	Args any
	Err  error
}
