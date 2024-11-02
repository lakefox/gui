package events

import (
	"fmt"
	adapter "gui/adapters"
	"gui/cstyle"
	"gui/element"
	"strings"
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
	State   *map[string]element.State
	Adapter *adapter.Adapter
	History *map[string]element.EventList
	CSS     *cstyle.CSS
	Focus   Focus
}

type Focus struct {
	Selected            int
	LastClickWasFocused bool
	Nodes               []*element.Node
	SoftFocused         *element.Node
}

// this could take the real document and apply the events calculated to it
func (m *Monitor) RunEvents(n *element.Node) {

}

// It does make sense that events would be attached to the state and not the document (the fake document)

func (m *Monitor) GetEvents(data *EventData) {
	s := *m.State

	eventMap := map[string][]element.Event{}

	for k, self := range s {
		if data.Click {
			m.Focus.SoftFocused = nil
		}

		var isMouseOver, isFocused bool

		if m.Focus.SoftFocused != nil {
			isFocused = m.Focus.SoftFocused.Properties.Id == k
		} else {
			isFocused = false
		}

		evt := element.Event{}

		if eventMap[k] == nil {
			eventMap[k] = []element.Event{}
		}

		insideX := (self.X < float32(data.Position[0]) && self.X+self.Width > float32(data.Position[0]))
		insideY := (self.Y < float32(data.Position[1]) && self.Y+self.Height > float32(data.Position[1]))
		inside := (insideX && insideY)

		arrowScroll := 0

		if isFocused {
			// This allows for arrow scrolling when off the element and typing
			if m.Focus.SoftFocused != nil {
				if data.Key == 265 {
					// up
					arrowScroll += 20
				} else if data.Key == 264 {
					// Down
					arrowScroll -= 20
				}
			}
			// Get the keycode of the pressed key
			if data.Key != 0 {
				if self.ContentEditable {
					// Sync the innertext and value but idk
					// !ISSUE: This may not work
					ProcessKeyEvent(self, int(data.Key))

					evt.Value = self.Value
				}
			}

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
		}

		if inside || isFocused {
			// Mouse is over element
			isMouseOver = true

			if evt.AddClasses == nil {
				evt.AddClasses = []string{}
			}

			evt.AddClasses = append(evt.AddClasses, ":hover")

			if data.Click && !evt.MouseDown {
				evt.MouseDown = true
				evt.MouseUp = false
			}

			if !data.Click && !evt.MouseUp {
				evt.MouseUp = true
				evt.MouseDown = false
				evt.Click = false
			}

			if data.Click && !evt.Click {
				evt.Click = true

				// if n.TabIndex > -1 {
				// 	if m.Focus.Selected > -1 {
				// 		if m.Focus.Nodes[m.Focus.Selected].Properties.Id != k {
				// 			m.Focus.Nodes[m.Focus.Selected].Blur()
				// 			for i, v := range m.Focus.Nodes {
				// 				if v.Properties.Id == k {
				// 					m.Focus.Selected = i
				// 					m.Focus.LastClickWasFocused = true
				// 					break
				// 				}
				// 			}
				// 		} else {
				// 			m.Focus.LastClickWasFocused = true
				// 		}
				// 	} else {
				// 		selectedIndex := -1
				// 		for i, v := range m.Focus.Nodes {
				// 			if v.Properties.Id == k {
				// 				selectedIndex = i
				// 			}
				// 		}
				// 		if selectedIndex == -1 {
				// 			if n.TabIndex == 9999999 {
				// 				// Add the last digits of the properties.id to make the elements sort in order
				// 				numStr := strings.TrimFunc(k, func(r rune) bool {
				// 					return !unicode.IsDigit(r) // Remove non-digit characters
				// 				})
				// 				prid, _ := strconv.Atoi(numStr)
				// 				n.TabIndex += prid
				// 			}
				// 			m.Focus.Nodes = append([]*element.Node{n}, m.Focus.Nodes...)
				// 			sort.Slice(m.Focus.Nodes, func(i, j int) bool {
				// 				return m.Focus.Nodes[i].TabIndex < m.Focus.Nodes[j].TabIndex // Ascending order by TabIndex
				// 			})
				// 			for i, v := range m.Focus.Nodes {
				// 				if v.Properties.Id == k {
				// 					selectedIndex = i
				// 				}
				// 			}
				// 		}

				// 		m.Focus.Selected = selectedIndex
				// 		m.Focus.LastClickWasFocused = true
				// 	}

				// } else if m.Focus.Selected > -1 {
				// 	if m.Focus.Nodes[m.Focus.Selected].Properties.Id != k && !m.Focus.LastClickWasFocused {
				// 		m.Focus.Nodes[m.Focus.Selected].Blur()
				// 		m.Focus.Selected = -1
				// 	}
				// }

				// Regardless set soft focus to trigger events to the selected element: when non is set default body???
				// if m.Focus.SoftFocused == nil {
				// 	m.Focus.SoftFocused = n
				// }
			}

			if data.Context {
				evt.ContextMenu = true
			}

			// el.ScrollY = 0
			// if (data.Scroll != 0 && (inside)) || arrowScroll != 0 {
			// 	// !TODO: for now just emit a event, will have to add el.scrollX
			// 	data.Scroll += arrowScroll
			// 	styledEl, _ := m.CSS.GetStyles(n)

			// 	// !TODO: Add scrolling for dragging over the scroll bar
			// 	// + the dragging part will be hard as events has no context of the scrollbars

			// 	if hasAutoOrScroll(styledEl) {
			// 		n.ScrollTop = int(n.ScrollTop + (-data.Scroll))
			// 		if n.ScrollTop > n.ScrollHeight {
			// 			n.ScrollTop = n.ScrollHeight
			// 		}

			// 		if n.ScrollTop <= 0 {
			// 			n.ScrollTop = 0
			// 		}

			// 		if n.OnScroll != nil {
			// 			n.OnScroll(evt)
			// 		}

			// 		data.Scroll = 0
			// 		eventList = append(eventList, "scroll")
			// 	}

			// }

			if !evt.MouseEnter {
				evt.MouseEnter = true
				evt.MouseOver = true
				evt.MouseLeave = false

				// Let the adapter know the cursor has changed
				// m.Adapter.DispatchEvent(element.Event{
				// 	Name: "cursor",
				// 	Data: self.Cursor,
				// })
			}

			if evt.X != int(data.Position[0]) && evt.Y != int(data.Position[1]) {
				evt.X = int(data.Position[0])
				evt.Y = int(data.Position[1])

			}
		} else {
			isMouseOver = false
			if evt.RemoveClasses == nil {
				evt.RemoveClasses = []string{}
			}

			evt.RemoveClasses = append(evt.RemoveClasses, ":hover")
		}

		if !isMouseOver && !evt.MouseLeave {
			evt.MouseEnter = false
			evt.MouseOver = false
			evt.MouseLeave = true
			// n.Properties.Hover = false
		}

		eventMap[k] = append(eventMap[k], evt)
	}

	fmt.Println(eventMap["ROOT"])
}

// ProcessKeyEvent processes key events for text entry.
func ProcessKeyEvent(self element.State, key int) {
	// Handle key events for text entry
	switch key {
	case 8:
		// Backspace: remove the last character
		if len(self.Value) > 0 {
			self.Value = self.Value[:len(self.Value)-1]
			// n.InnerText = n.InnerText[:len(n.InnerText)-1]
		}

	case 65:
		// !TODO: ctrl a
		// // Select All: set the entire text as selected
		// if key == 17 || key == 345 {
		// 	n.Properties.Selected = []float32{0, float32(len(n.Value))}
		// } else {
		// 	// Otherwise, append 'A' to the text
		// 	n.Value += "A"
		// }

	case 67:
		// Copy: copy the selected text (in this case, print it)
		// if key == 17 || key == 345 {
		// 	fmt.Println("Copy:", n.Value)
		// } else {
		// 	// Otherwise, append 'C' to the text
		// 	n.Value += "C"
		// }

	case 86:
		// Paste: paste the copied text (in this case, set it to "Pasted")
		// if key == 17 || key == 345 {
		// 	n.Value = "Pasted"
		// } else {
		// 	// Otherwise, append 'V' to the text
		// 	n.Value += "V"
		// }

	default:
		// Record other key presses
		self.Value += string(rune(key))
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
