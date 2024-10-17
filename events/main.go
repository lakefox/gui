package events

import (
	"fmt"
	adapter "gui/adapters"
	"gui/cstyle"
	"gui/element"
	"sort"
	"strconv"
	"strings"
	"unicode"

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
	Focus     Focus
}

type Focus struct {
	Selected            int
	LastClickWasFocused bool
	Nodes               []*element.Node
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
	m.Focus.LastClickWasFocused = false
	m.CalcEvents(m.Document, data)
}

func (m *Monitor) CalcEvents(n *element.Node, data *EventData) {
	// loop through state to build events, then use multithreading to complete
	// map
	for _, v := range n.Children {
		m.CalcEvents(v, data)
	}

	if n.TabIndex > -1 {
		if m.Focus.Nodes != nil {
			selectedIndex := -1
			for i, v := range m.Focus.Nodes {
				if v.Properties.Id == n.Properties.Id {
					selectedIndex = i
				}
			}
			if selectedIndex == -1 {
				if n.TabIndex == 9999999 {
					// Add the last digits of the properties.id to make the elements sort in order
					numStr := strings.TrimFunc(n.Properties.Id, func(r rune) bool {
						return !unicode.IsDigit(r) // Remove non-digit characters
					})
					prid, _ := strconv.Atoi(numStr)
					n.TabIndex += prid
				}
				m.Focus.Nodes = append([]*element.Node{n}, m.Focus.Nodes...)
			}
		} else {
			m.Focus = Focus{}
			m.Focus.Selected = -1
			m.Focus.Nodes = append(m.Focus.Nodes, n)
		}
	}
	sort.Slice(m.Focus.Nodes, func(i, j int) bool {
		return m.Focus.Nodes[i].TabIndex < m.Focus.Nodes[j].TabIndex // Ascending order by TabIndex
	})

	mHistory := *m.History
	eventList := []string{}
	evt := mHistory[n.Properties.Id].Event

	s := *m.State
	self := s[n.Properties.Id]

	if evt.Target == nil {
		if len(data.Position) < 2 {
			data.Position = []int{0, 0}
		}
		evt = element.Event{
			X:          data.Position[0],
			Y:          data.Position[1],
			MouseUp:    true,
			MouseLeave: true,
			Target:     n,
		}
		(*m.History)[n.Properties.Id] = element.EventList{
			Event: evt,
			List:  []string{},
		}
	}

	var isMouseOver bool

	if self.X < float32(data.Position[0]) && self.X+self.Width > float32(data.Position[0]) {
		if self.Y < float32(data.Position[1]) && self.Y+self.Height > float32(data.Position[1]) {
			// Mouse is over element
			isMouseOver = true
			if !slices.Contains(n.ClassList.Classes, ":hover") {
				n.ClassList.Add(":hover")
			}

			if data.Click && !evt.MouseDown {
				evt.MouseDown = true
				evt.MouseUp = false
				if n.OnMouseDown != nil {
					n.OnMouseDown(evt)
				}
				eventList = append(eventList, "mousedown")
			}

			if !data.Click && !evt.MouseUp {
				evt.MouseUp = true
				evt.MouseDown = false
				evt.Click = false
				if n.OnMouseUp != nil {
					n.OnMouseUp(evt)
				}
				eventList = append(eventList, "mouseup")
			}

			if data.Click && !evt.Click {
				evt.Click = true

				if n.TabIndex > -1 {
					if m.Focus.Selected > -1 {
						if m.Focus.Nodes[m.Focus.Selected].Properties.Id != n.Properties.Id {
							m.Focus.Nodes[m.Focus.Selected].Blur()
							for i, v := range m.Focus.Nodes {
								if v.Properties.Id == n.Properties.Id {
									m.Focus.Selected = i
									n.Focus()
									m.Focus.LastClickWasFocused = true
									break
								}
							}
						} else {
							m.Focus.LastClickWasFocused = true
						}
					} else {
						selectedIndex := -1
						for i, v := range m.Focus.Nodes {
							if v.Properties.Id == n.Properties.Id {
								selectedIndex = i
							}
						}
						if selectedIndex == -1 {
							if n.TabIndex == 9999999 {
								// Add the last digits of the properties.id to make the elements sort in order
								numStr := strings.TrimFunc(n.Properties.Id, func(r rune) bool {
									return !unicode.IsDigit(r) // Remove non-digit characters
								})
								prid, _ := strconv.Atoi(numStr)
								n.TabIndex += prid
							}
							m.Focus.Nodes = append([]*element.Node{n}, m.Focus.Nodes...)
							sort.Slice(m.Focus.Nodes, func(i, j int) bool {
								return m.Focus.Nodes[i].TabIndex < m.Focus.Nodes[j].TabIndex // Ascending order by TabIndex
							})
							for i, v := range m.Focus.Nodes {
								if v.Properties.Id == n.Properties.Id {
									selectedIndex = i
								}
							}
						}

						m.Focus.Selected = selectedIndex
						n.Focus()
						m.Focus.LastClickWasFocused = true
					}

				} else if m.Focus.Selected > -1 {
					if m.Focus.Nodes[m.Focus.Selected].Properties.Id != n.Properties.Id && !m.Focus.LastClickWasFocused {
						m.Focus.Nodes[m.Focus.Selected].Blur()
						m.Focus.Selected = -1
					}
				}
				if n.OnClick != nil {
					n.OnClick(evt)
				}
				eventList = append(eventList, "click")
			}

			if data.Context {
				evt.ContextMenu = true
				if n.OnContextMenu != nil {
					n.OnContextMenu(evt)
				}
				eventList = append(eventList, "contextmenu")
			}

			// el.ScrollY = 0
			if data.Scroll != 0 {
				// !TODO: for now just emit a event, will have to add el.scrollX

				styledEl, _ := m.CSS.GetStyles(n)

				// !TODO: Add scrolling for dragging over the scroll bar and arrow keys if it is focused
				// + Working on the focus part, the dragging part will be hard as events has no context of the scrollbars

				if hasAutoOrScroll(styledEl) {
					n.ScrollTop = int(n.ScrollTop + (-data.Scroll))
					if n.ScrollTop > n.ScrollHeight {
						n.ScrollTop = n.ScrollHeight
					}

					if n.ScrollTop <= 0 {
						n.ScrollTop = 0
					}

					if n.OnScroll != nil {
						n.OnScroll(evt)
					}

					data.Scroll = 0

					eventList = append(eventList, "scroll")
				}

			}

			if !evt.MouseEnter {
				evt.MouseEnter = true
				evt.MouseOver = true
				evt.MouseLeave = false
				if n.OnMouseEnter != nil {
					n.OnMouseEnter(evt)
				}
				if n.OnMouseOver != nil {
					n.OnMouseEnter(evt)
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
				if n.OnMouseMove != nil {
					n.OnMouseMove(evt)
				}
				eventList = append(eventList, "mousemove")
			}

			// Get the keycode of the pressed key
			// issue: need to only add the text data and events to focused elements and only
			// 		  one at a time
			if data.Key != 0 {
				if n.ContentEditable {
					n.Value = n.InnerText
					ProcessKeyEvent(n, int(data.Key))

					n.InnerText = n.Value
					eventList = append(eventList, "keypress")
				}
			}

			// !ISSUE: Tab fires multiple events
			if data.Key == 258 && data.KeyState && !m.Focus.LastClickWasFocused {
				// Tab
				mfsLen := len(m.Focus.Nodes)
				if mfsLen > 0 {
					store := m.Focus.Selected
					m.Focus.Selected += 1
					if m.Focus.Selected >= mfsLen {
						m.Focus.Selected = 0
					}
					if store != m.Focus.Selected {
						if store > -1 {
							m.Focus.Nodes[store].Blur()
						}
						m.Focus.Nodes[m.Focus.Selected].Focus()
						m.Focus.LastClickWasFocused = true
					}
				}

			}

		} else {
			isMouseOver = false
			if slices.Contains(n.ClassList.Classes, ":hover") {
				n.ClassList.Remove(":hover")
			}
		}
	} else {
		isMouseOver = false
		if slices.Contains(n.ClassList.Classes, ":hover") {
			n.ClassList.Remove(":hover")
		}
	}

	// fmt.Println(isMouseOver)

	if !isMouseOver && !evt.MouseLeave {
		evt.MouseEnter = false
		evt.MouseOver = false
		evt.MouseLeave = true
		n.Properties.Hover = false
		if n.OnMouseLeave != nil {
			n.OnMouseLeave(evt)
		}
		eventList = append(eventList, "mouseleave")
	}

	if len(n.Properties.Events) > 0 {
		eventList = append(eventList, n.Properties.Events...)
	}

	(*m.History)[n.Properties.Id] = element.EventList{
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

func hasAutoOrScroll(styledEl map[string]string) bool {
	overflowKeys := []string{"overflow", "overflow-x", "overflow-y"}
	for _, key := range overflowKeys {
		if value, exists := styledEl[key]; exists {
			values := strings.Fields(value) // Splits the value by spaces
			for _, v := range values {
				if v == "auto" || v == "scroll" {
					return true
				}
			}
		}
	}
	return false
}
