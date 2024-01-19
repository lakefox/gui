package events

import (
	"fmt"
	"gui/element"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func GetEvents(el *element.Node, prevEvents *map[string]element.Event) *map[string]element.Event {
	mp := rl.GetMousePosition()
	// Mouse over
	// fmt.Println(len(*prevEvents))
	loop(el, mp, prevEvents)
	return prevEvents
}

func loop(el *element.Node, mp rl.Vector2, eventTracker *map[string]element.Event) {
	et := *eventTracker
	eventList := []string{}
	evt := et[el.Properties.Id]

	if evt.Target.Properties.Id == "" {
		et[el.Properties.Id] = element.Event{
			X:          int(mp.X),
			Y:          int(mp.Y),
			MouseUp:    true,
			MouseLeave: true,
			Target:     *el,
		}

		evt = et[el.Properties.Id]
	}

	var isMouseOver bool

	if el.Properties.X < mp.X && el.Properties.X+el.Properties.Width > mp.X {
		if el.Properties.Y < mp.Y && el.Properties.Y+el.Properties.Height > mp.Y {
			// Mouse is over element
			isMouseOver = true

			fmt.Println(rl.GetMouseWheelMove())

			if rl.IsMouseButtonDown(rl.MouseLeftButton) && !evt.MouseDown {
				evt.MouseDown = true
				evt.MouseUp = false
				if el.OnMouseDown != nil {
					el.OnMouseDown(evt)
				}
				eventList = append(eventList, "mousedown")
			}

			if !rl.IsMouseButtonDown(rl.MouseLeftButton) && !evt.MouseUp {
				evt.MouseUp = true
				evt.MouseDown = false
				if el.OnMouseUp != nil {
					el.OnMouseUp(evt)
				}
				eventList = append(eventList, "mouseup")
			}

			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				evt.Click = true
				if el.OnClick != nil {
					el.OnClick(evt)
				}
				eventList = append(eventList, "click")
			}

			if rl.IsMouseButtonPressed(rl.MouseRightButton) {
				evt.ContextMenu = true
				if el.OnContextMenu != nil {
					el.OnContextMenu(evt)
				}
				eventList = append(eventList, "contextmenu")
			}

			if !evt.MouseEnter {
				evt.MouseEnter = true
				evt.MouseOver = true
				evt.MouseLeave = false
				if el.OnMouseEnter != nil {
					el.OnMouseEnter(evt)
				}
				if el.OnMouseOver != nil {
					el.OnMouseEnter(evt)
				}
				eventList = append(eventList, "mouseenter")
				eventList = append(eventList, "mouseover")
			}

			if evt.X != int(mp.X) && evt.Y != int(mp.Y) {
				evt.X = int(mp.X)
				evt.Y = int(mp.Y)
				if el.OnMouseMove != nil {
					el.OnMouseMove(evt)
				}
				eventList = append(eventList, "mousemove")
			}

			// Get the keycode of the pressed key
			keyPressed := rl.GetKeyPressed()
			if keyPressed != 0 {
				fmt.Printf("Key pressed: %c (%d)\n", keyPressed, keyPressed)
			}

		} else {
			isMouseOver = false
		}
	} else {
		isMouseOver = false
	}

	// fmt.Println(isMouseOver)

	if !isMouseOver && !evt.MouseLeave {
		evt.MouseEnter = false
		evt.MouseOver = false
		evt.MouseLeave = true
		if el.OnMouseLeave != nil {
			el.OnMouseLeave(evt)
		}
		eventList = append(eventList, "mouseleave")
	}

	if len(eventList) > 0 {
		for _, v := range eventList {
			if len(el.Properties.EventListeners[v]) > 0 {
				for _, handler := range el.Properties.EventListeners[v] {
					handler(evt)
				}
			}
		}

	}

	et[el.Properties.Id] = evt

	eventTracker = &et
	for _, v := range el.Children {
		loop(&v, mp, eventTracker)
	}
}
