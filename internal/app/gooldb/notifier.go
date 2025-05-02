package gooldb

type NotificationLevel int

const (
	Info NotificationLevel = iota
	Warning
	Panic
)

type Notifier interface {
	Notify(level NotificationLevel, args ...any)
}
