package main

import (
	"fmt"
	"time"

	"github.com/1071496910/event-controller/controller"
)

type FakeSucker struct{}

func (fs *FakeSucker) GetSucker() chan *controller.Event {
	ch := make(chan *controller.Event, 1024)
	go func() {
		for i := 0; ; i++ {
			ch <- &controller.Event{
				OPType: i % 3,
				Obj:    i,
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return ch
}

func main() {
	// func NewController(sucker controller.EventSucker, addFunc controller.EventDealFunc, updateFunc controller.EventDealFunc, delFunc controller.EventDealFunc, eventQueueLen
	ctl := controller.NewController(&FakeSucker{},
		func(i interface{}) {
			fmt.Println("In add func", i)
		},
		func(i interface{}) {
			fmt.Println("In update func", i)
		},
		func(i interface{}) {
			fmt.Println("In del func", i)
		},
		1024)
	go func() {
		time.Sleep(10 * time.Second)
		ctl.Stop()
	}()
	ctl.Run()

}
