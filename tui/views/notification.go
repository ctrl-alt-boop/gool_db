package views

import (
	"fmt"
	"time"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/jesseduffield/gocui"
)

const (
	InfoBackground    = gocui.ColorDefault
	WarningBackground = gocui.ColorYellow
	PanicBackground   = gocui.ColorRed
)

const NotificationViewPrefix = "notification_"

// const maxShowNotifications = 3

type NotificationView struct {
	*gocui.View
	Level gooldb.NotificationLevel
	Name  string
}

type NotificationHandler struct {
	gui           *gocui.Gui
	notifications []*NotificationView
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notifications: make([]*NotificationView, 0),
	}
}

func (n *NotificationHandler) NewNotification(level gooldb.NotificationLevel, args ...any) {
	maxX, maxY := n.gui.Size()
	notificationName := fmt.Sprintf("%s%d", NotificationViewPrefix, len(n.notifications))

	view, err := n.gui.SetView(notificationName, maxX/3, 0, (maxX/3)*2, maxY-4, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			//panic(err)
		}
		switch level {
		case gooldb.Info:
			view.BgColor = InfoBackground
		case gooldb.Warning:
			view.BgColor = WarningBackground
		case gooldb.Panic:
			view.BgColor = PanicBackground
		}
		view.Wrap = true
		view.Frame = true
		view.Editable = false
		fmt.Fprint(view, args...)
		notificationView := &NotificationView{
			View:  view,
			Level: level,
			Name:  notificationName,
		}
		n.notifications = append(n.notifications, notificationView)
		time.AfterFunc(5*time.Second, func() {
			n.removeTopNotification(notificationName)
		})
	}
}

func (n *NotificationHandler) removeTopNotification(name string) {

	n.gui.DeleteView(name)
	if len(n.notifications) == 0 {
		return
	}
	n.notifications = n.notifications[1:]
}
