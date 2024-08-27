package adapter

import (
	"gui/element"
	"gui/library"
)

type Adapter struct {
	Init    func(width int, height int)
	Render  func(state []element.State)
	Load    func(state []element.State)
	events  map[string][]func(element.Event)
	Library *library.Shelf
	Options Options
}

type Options struct {
	RenderText     bool
	RenderElements bool
	RenderBorders  bool
}

func (a *Adapter) AddEventListener(name string, callback func(element.Event)) {
	if a.events == nil {
		a.events = map[string][]func(element.Event){}
	}
	a.events[name] = append(a.events[name], callback)
}

func (a *Adapter) DispatchEvent(event element.Event) {
	// fmt.Println("here", a.events)
	if a.events != nil {
		evts := a.events[event.Name]
		for _, v := range evts {
			v(event)
		}
	}
}
