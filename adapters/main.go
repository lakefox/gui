package adapter

import (
	"gui/element"
	"gui/library"
)

// !TODO: For things like fonts the file system interaction works fine, but if you want to run grim on a pi pico it doesn't have the file system int like
// + computers have
// + Option 1: make a file system adapter where authors can control how data is read
// + Option 2: at build time fetch all files needed and bundle them

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
