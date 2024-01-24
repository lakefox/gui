package events

import (
	"fmt"
	"gui/element"
	"gui/utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RLData struct {
	MP rl.Vector2
	LB bool
	RB bool
	WD float32
	KP int32
}

func GetEvents(el *element.Node, prevEvents *map[string]element.EventList) *map[string]element.EventList {
	data := RLData{
		MP: rl.GetMousePosition(),
		LB: rl.IsMouseButtonDown(rl.MouseLeftButton),
		RB: rl.IsMouseButtonPressed(rl.MouseRightButton),
		WD: rl.GetMouseWheelMove(),
		KP: rl.GetKeyPressed(),
	}
	// Mouse over
	// fmt.Println(len(*prevEvents))
	loop(el, data, prevEvents)
	return prevEvents
}

func RunEvents(events *map[string]element.EventList) {
	for _, evt := range *events {
		if len(evt.List) > 0 {
			for _, v := range evt.List {
				if len(evt.Event.Target.Properties.EventListeners[v]) > 0 {
					for _, handler := range evt.Event.Target.Properties.EventListeners[v] {
						handler(evt.Event)
					}
				}
			}
		}
	}

}

func loop(el *element.Node, data RLData, eventTracker *map[string]element.EventList) *element.Node {
	et := *eventTracker
	eventList := []string{}
	evt := et[el.Properties.Id].Event

	if evt.Target.Properties.Id == "" {
		et[el.Properties.Id] = element.EventList{
			Event: element.Event{
				X:          int(data.MP.X),
				Y:          int(data.MP.Y),
				MouseUp:    true,
				MouseLeave: true,
				Target:     *el,
			},
			List: []string{},
		}

		evt = et[el.Properties.Id].Event
	}

	var isMouseOver bool

	if el.Properties.X < data.MP.X && el.Properties.X+el.Properties.Width > data.MP.X {
		if el.Properties.Y < data.MP.Y && el.Properties.Y+el.Properties.Height > data.MP.Y {
			// Mouse is over element
			isMouseOver = true

			if data.LB && !evt.MouseDown {
				evt.MouseDown = true
				evt.MouseUp = false
				if el.OnMouseDown != nil {
					el.OnMouseDown(evt)
				}
				eventList = append(eventList, "mousedown")
			}

			if !data.LB && !evt.MouseUp {
				evt.MouseUp = true
				evt.MouseDown = false
				evt.Click = false
				if el.OnMouseUp != nil {
					el.OnMouseUp(evt)
				}
				eventList = append(eventList, "mouseup")
			}

			if data.LB && !evt.Click {
				evt.Click = true
				if el.OnClick != nil {
					el.OnClick(evt)
				}
				eventList = append(eventList, "click")
			}

			if data.RB {
				evt.ContextMenu = true
				if el.OnContextMenu != nil {
					el.OnContextMenu(evt)
				}
				eventList = append(eventList, "contextmenu")
			}

			if data.WD != 0 {
				// fmt.Println(data.WD)
				// for now just emit a event, will have to add el.scrollX
				evt.Target.ScrollY = utils.Max(evt.Target.ScrollY+(-data.WD), 0)

				if el.OnScroll != nil {
					el.OnScroll(evt)
				}
				eventList = append(eventList, "scroll")
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

			if evt.X != int(data.MP.X) && evt.Y != int(data.MP.Y) {
				evt.X = int(data.MP.X)
				evt.Y = int(data.MP.Y)
				if el.OnMouseMove != nil {
					el.OnMouseMove(evt)
				}
				eventList = append(eventList, "mousemove")
			}

			// Get the keycode of the pressed key
			// issue: need to only add the text data and events to focused elements and only
			// 		  one at a time
			if data.KP != 0 {
				if el.Properties.Editable {
					el.Value = el.InnerText
					ProcessKeyEvent(el, int(data.KP))
					fmt.Println(el.Properties.Id, el.Value)

					el.InnerText = el.Value
					eventList = append(eventList, "keypress")
				}

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

	et[el.Properties.Id] = element.EventList{
		Event: evt,
		List:  eventList,
	}

	eventTracker = &et
	for i, v := range el.Children {
		el.Children[i] = *loop(&v, data, eventTracker)
	}
	return el
}

// ProcessKeyEvent processes key events for text entry.
func ProcessKeyEvent(n *element.Node, key int) {
	// Handle key events for text entry
	switch key {
	case rl.KeyBackspace:
		// Backspace: remove the last character
		if len(n.Value) > 0 {
			n.Value = n.Value[:len(n.Value)-1]
		}

	case rl.KeyA:
		// Select All: set the entire text as selected
		if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
			n.Properties.Selected = []float32{0, float32(len(n.Value))}
		} else {
			// Otherwise, append 'A' to the text
			n.Value += "A"
		}

	case rl.KeyC:
		// Copy: copy the selected text (in this case, print it)
		if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
			fmt.Println("Copy:", n.Value)
		} else {
			// Otherwise, append 'C' to the text
			n.Value += "C"
		}

	case rl.KeyV:
		// Paste: paste the copied text (in this case, set it to "Pasted")
		if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
			n.Value = "Pasted"
		} else {
			// Otherwise, append 'V' to the text
			n.Value += "V"
		}

	default:
		// Record other key presses
		n.Value += string(rune(key))
	}
}
