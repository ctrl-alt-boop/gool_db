package gooldb

type State int

const (
	Database State = iota
	Table
)

var GoolState State = Database
