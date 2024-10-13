package events

import (
	"fmt"
	adapter "gui/adapters"
	"gui/cstyle"
	"gui/element"

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

type Monitor struct {
	Document  *element.Node
	State     *map[string]element.State
	EventList *element.EventList
	Adapter   *adapter.Adapter
	History   *map[string]element.EventList
	CSS       *cstyle.CSS
}

func (m *Monitor) RunEvents() bool {
	eventRan := false
	for _, evt := range *m.History {
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

func (m *Monitor) GetEvents(data *EventData) {
	m.CalcEvents(m.Document, data)
}

func (m *Monitor) CalcEvents(el *element.Node, data *EventData) {
	// loop through state to build events, then use multithreading to complete
	// map
	for _, v := range el.Children {
		m.CalcEvents(v, data)
	}

	mHistory := *m.History
	eventList := []string{}
	evt := mHistory[el.Properties.Id].Event

	s := *m.State
	self := s[el.Properties.Id]

	if evt.Target == nil {
		evt = element.Event{
			X:          data.Position[0],
			Y:          data.Position[1],
			MouseUp:    true,
			MouseLeave: true,
			Target:     el,
		}
		(*m.History)[el.Properties.Id] = element.EventList{
			Event: evt,
			List:  []string{},
		}
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

			// el.ScrollY = 0
			if data.Scroll != 0 {
				// !TODO: for now just emit a event, will have to add el.scrollX

				styledEl, _ := m.CSS.GetStyles(el)

				if styledEl["overflow"] != "" || styledEl["overflow-x"] != "" || styledEl["overflow-y"] != "" {
					el.ScrollTop = int(el.ScrollTop + (-data.Scroll))
					if el.ScrollTop > el.ScrollHeight-int(self.Height) {
						el.ScrollTop = el.ScrollHeight - int(self.Height)
					}

					if el.ScrollTop <= 0 {
						el.ScrollTop = 0
					}

					if el.OnScroll != nil {
						el.OnScroll(evt)
					}

					data.Scroll = 0

					eventList = append(eventList, "scroll")
				}

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
				m.Adapter.DispatchEvent(element.Event{
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

	if len(el.Properties.Events) > 0 {
		eventList = append(eventList, el.Properties.Events...)
	}

	(*m.History)[el.Properties.Id] = element.EventList{
		Event: evt,
		List:  eventList,
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
