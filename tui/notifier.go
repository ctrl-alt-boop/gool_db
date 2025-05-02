package tui

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

type Notifier struct {
	onNotification func(gooldb.NotificationLevel, ...any)
}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) Notify(level gooldb.NotificationLevel, args ...any) {
	n.onNotification(level, args...)
}

func (n *Notifier) SetOnNotification(f func(gooldb.NotificationLevel, ...any)) {
	n.onNotification = f
}
