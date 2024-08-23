package events

import (
	"fmt"
	adapter "gui/adapters"
	"gui/element"
	"gui/utils"

	"slices"
)

type EventData struct {
	Position []int
	Click    bool
	Context  bool
	Scroll   int
	Key      int
	KeyState bool
}

// func GetEvents(el *element.Node, state *map[string]element.State, prevEvents *map[string]element.EventList) *map[string]element.EventList {
// 	data := RLData{
// 		MP: rl.GetMousePosition(),
// 		LB: rl.IsMouseButtonDown(rl.MouseLeftButton),
// 		RB: rl.IsMouseButtonPressed(rl.MouseRightButton),
// 		WD: rl.GetMouseWheelMove(),
// 		KP: rl.GetKeyPressed(),
// 	}
// 	// fmt.Println(data.WD)
// 	// Mouse over
// 	// fmt.Println(len(*prevEvents))
// 	loop(el, state, data, prevEvents)
// 	return prevEvents
// }

func RunEvents(events *map[string]element.EventList) bool {
	eventRan := false
	for _, evt := range *events {
		if len(evt.List) > 0 {
			for _, v := range evt.List {
				if len(evt.Event.Target.Properties.EventListeners[v]) > 0 {
					for _, handler := range evt.Event.Target.Properties.EventListeners[v] {
						handler(evt.Event)
						eventRan = true
					}
				}
			}
		}
	}
	return eventRan
}

func GetEvents(el *element.Node, state *map[string]element.State, data EventData, eventTracker *map[string]element.EventList, a *adapter.Adapter) {
	// loop through state to build events, then use multithreading to complete
	// map
	et := *eventTracker
	eventList := []string{}
	evt := et[el.Properties.Id].Event

	s := *state
	self := s[el.Properties.Id]

	if evt.Target == nil {
		et[el.Properties.Id] = element.EventList{
			Event: element.Event{
				X:          data.Position[0],
				Y:          data.Position[1],
				MouseUp:    true,
				MouseLeave: true,
				Target:     el,
			},
			List: []string{},
		}

		evt = et[el.Properties.Id].Event
	}

	var isMouseOver bool

	if self.X < float32(data.Position[0]) && self.X+self.Width > float32(data.Position[0]) {
		if self.Y < float32(data.Position[1]) && self.Y+self.Height > float32(data.Position[1]) {
			// Mouse is over element
			isMouseOver = true
			if !slices.Contains(el.ClassList.Classes, ":hover") {
				el.ClassList.Add(":hover")
			}

			if data.Click && !evt.MouseDown {
				evt.MouseDown = true
				evt.MouseUp = false
				if el.OnMouseDown != nil {
					el.OnMouseDown(evt)
				}
				eventList = append(eventList, "mousedown")
			}

			if !data.Click && !evt.MouseUp {
				evt.MouseUp = true
				evt.MouseDown = false
				evt.Click = false
				if el.OnMouseUp != nil {
					el.OnMouseUp(evt)
				}
				eventList = append(eventList, "mouseup")
			}

			if data.Click && !evt.Click {
				evt.Click = true
				if el.OnClick != nil {
					el.OnClick(evt)
				}
				eventList = append(eventList, "click")
			}

			if data.Context {
				evt.ContextMenu = true
				if el.OnContextMenu != nil {
					el.OnContextMenu(evt)
				}
				eventList = append(eventList, "contextmenu")
			}

			if data.Scroll != 0 {
				// fmt.Println(data.WD)
				// for now just emit a event, will have to add el.scrollX
				evt.Target.ScrollY = int(utils.Max(float32(evt.Target.ScrollY+(-data.Scroll)), 0.0))

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

				// Let the adapter know the cursor has changed
				a.DispatchEvent(element.Event{
					Name: "cursor",
					Data: self.Cursor,
				})
			}

			if evt.X != int(data.Position[0]) && evt.Y != int(data.Position[1]) {
				evt.X = int(data.Position[0])
				evt.Y = int(data.Position[1])
				if el.OnMouseMove != nil {
					el.OnMouseMove(evt)
				}
				eventList = append(eventList, "mousemove")
			}

			// Get the keycode of the pressed key
			// issue: need to only add the text data and events to focused elements and only
			// 		  one at a time
			if data.Key != 0 {
				if el.Properties.Editable {
					el.Value = el.InnerText
					ProcessKeyEvent(el, int(data.Key))
					fmt.Println(el.Properties.Id, el.Value)

					el.InnerText = el.Value
					eventList = append(eventList, "keypress")
				}

			}

		} else {
			isMouseOver = false
			if slices.Contains(el.ClassList.Classes, ":hover") {
				el.ClassList.Remove(":hover")
			}
		}
	} else {
		isMouseOver = false
		if slices.Contains(el.ClassList.Classes, ":hover") {
			el.ClassList.Remove(":hover")
		}
	}

	// fmt.Println(isMouseOver)

	if !isMouseOver && !evt.MouseLeave {
		evt.MouseEnter = false
		evt.MouseOver = false
		evt.MouseLeave = true
		el.Properties.Hover = false
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
	for _, v := range el.Children {
		GetEvents(v, state, data, eventTracker, a)
	}
}

// ProcessKeyEvent processes key events for text entry.
func ProcessKeyEvent(n *element.Node, key int) {
	// Handle key events for text entry
	switch key {
	case 8:
		// Backspace: remove the last character
		if len(n.Value) > 0 {
			n.Value = n.Value[:len(n.Value)-1]
			n.InnerText = n.InnerText[:len(n.InnerText)-1]
		}

	case 65:
		// Select All: set the entire text as selected
		if key == 17 || key == 345 {
			n.Properties.Selected = []float32{0, float32(len(n.Value))}
		} else {
			// Otherwise, append 'A' to the text
			n.Value += "A"
		}

	case 67:
		// Copy: copy the selected text (in this case, print it)
		if key == 17 || key == 345 {
			fmt.Println("Copy:", n.Value)
		} else {
			// Otherwise, append 'C' to the text
			n.Value += "C"
		}

	case 86:
		// Paste: paste the copied text (in this case, set it to "Pasted")
		if key == 17 || key == 345 {
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
