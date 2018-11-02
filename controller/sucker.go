package controller

type EventSucker interface {
	GetSucker() chan *Event
}
