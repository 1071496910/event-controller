package controller

import (
	"fmt"
	"sync"
)

const (
	op_add int = iota
	op_update
	op_del
)

type Event struct {
	OPType int
	Obj    interface{}
}

type EventDealFunc func(e *Event)

type Controller struct {
	sucker EventSucker

	addFunc    EventDealFunc
	updateFunc EventDealFunc
	delFunc    EventDealFunc

	addChan    chan *Event
	updateChan chan *Event
	delChan    chan *Event
	stopChan   chan int
}

func NewController(sucker EventSucker, addFunc EventDealFunc, updateFunc EventDealFunc, delFunc EventDealFunc, eventQueueLen int) *Controller {
	return &Controller{
		sucker:     sucker,
		addFunc:    addFunc,
		updateFunc: updateFunc,
		delFunc:    delFunc,
		addChan:    make(chan *Event, eventQueueLen),
		updateChan: make(chan *Event, eventQueueLen),
		delChan:    make(chan *Event, eventQueueLen),
		stopChan:   make(chan int),
	}
}

func (c *Controller) Run() error {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		eventChan := c.sucker.GetSucker()
		if eventChan == nil {
			panic(fmt.Errorf("get sucker error"))
		}
	SuckLoop:
		for {
			select {
			case e := <-eventChan:
				switch e.OPType {
				case op_add:
					c.addChan <- e
				case op_update:
					c.updateChan <- e
				case op_del:
					c.delChan <- e
				}

			case _ = <-c.stopChan:
				break SuckLoop
			}
		}
	}()

	go func() {
		defer wg.Done()
	AddLoop:
		for {
			select {
			case e := <-c.addChan:
				c.addFunc(e)
			case _ = <-c.stopChan:
				break AddLoop
			}

		}
	}()
	go func() {
		defer wg.Done()
	UpdateLoop:
		for {
			select {
			case e := <-c.updateChan:
				c.updateFunc(e)
			case _ = <-c.stopChan:
				break UpdateLoop
			}
		}
	}()
	go func() {
		defer wg.Done()
	DelLoop:
		for {
			select {
			case e := <-c.delChan:
				c.delFunc(e)
			case _ = <-c.stopChan:
				break DelLoop
			}

		}
	}()
	wg.Wait()
	return nil
}

func (c *Controller) Stop() {
	close(c.stopChan)
}
