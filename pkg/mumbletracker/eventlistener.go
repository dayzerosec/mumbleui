package mumbletracker

import (
	"layeh.com/gumble/gumble"
	"sync"
)

func NewEventListener(OnJoin, OnLeave func(*gumble.User, []string)) eventListener {
	return eventListener{
		OnJoin:  OnJoin,
		OnLeave: OnLeave,
		lock:    &sync.Mutex{},
		users:   &map[uint32]*gumble.User{},
	}
}

type eventListener struct {
	OnJoin func(*gumble.User, []string)
	OnLeave func(*gumble.User, []string)
	lock  *sync.Mutex
	users *map[uint32]*gumble.User
}

func (l eventListener) userlist_locked() []string {
	var out []string
	for _, user := range *l.users {
		out = append(out, user.Name)
	}
	return out
}

func (l eventListener) OnConnect(e *gumble.ConnectEvent) {
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, user := range e.Client.Users {
		(*l.users)[user.UserID] = user
		l.OnJoin(user, l.userlist_locked())
	}
}
func (l eventListener) OnDisconnect(e *gumble.DisconnectEvent)   {}
func (l eventListener) OnTextMessage(e *gumble.TextMessageEvent) {}
func (l eventListener) OnUserChange(e *gumble.UserChangeEvent) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if e.Type.Has(gumble.UserChangeConnected) {
		for _, user := range e.Client.Users {
			if _, found := (*l.users)[user.UserID]; !found {
				(*l.users)[user.UserID] = user
				go l.OnJoin(user, l.userlist_locked())
			}
		}
	} else if e.Type.Has(gumble.UserChangeDisconnected) {
		for _, user := range *l.users {
			if e.Client.Users.Find(user.Name) == nil {
				delete(*l.users, user.UserID)
				go l.OnLeave(user, l.userlist_locked())
			}
		}
	}
}

func (l eventListener) OnChannelChange(e *gumble.ChannelChangeEvent)             {}
func (l eventListener) OnPermissionDenied(e *gumble.PermissionDeniedEvent)       {}
func (l eventListener) OnUserList(e *gumble.UserListEvent)                       {}
func (l eventListener) OnACL(e *gumble.ACLEvent)                                 {}
func (l eventListener) OnBanList(e *gumble.BanListEvent)                         {}
func (l eventListener) OnContextActionChange(e *gumble.ContextActionChangeEvent) {}
func (l eventListener) OnServerConfig(e *gumble.ServerConfigEvent)               {}
