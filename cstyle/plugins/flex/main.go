package flex

import (
	"gui/cstyle"
	"gui/element"
	"gui/utils"
	"slices"
	"strings"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Selector: func(n *element.Node) bool {
			styles := map[string]string{
				"display":         "flex",
				"justify-content": "*",
				"align-content":   "*",
				"align-items":     "*",
				"flex-wrap":       "*",
				"flex-direction":  "*",
			}
			matches := true
			for name, value := range styles {
				if n.Style[name] != value && !(value == "*") && n.Style[name] != "" {
					matches = false
				}
			}
			return matches
		},
		Level: 2,
		Handler: func(n *element.Node, state *map[string]element.State) {
			// !ISSUE: align-items is not impleamented
			s := *state
			self := s[n.Properties.Id]
			// Brief: justify does not align the bottom row correctly
			//        y axis also needs to be done
			verbs := strings.Split(n.Style["flex-direction"], "-")

			orderedNode := order(*n, state, n.Children, verbs[0], len(verbs) > 1, n.Style["flex-wrap"] == "wrap")

			// var i int

			colWidth := self.Width / float32(len(orderedNode))

			var xOffset, yOffset float32
			// if n.Style["justify-content"] == "space-evenly" {
			// 	b, _ := utils.ConvertToPixels(n.Children[i].Properties.Border.Width, n.Children[i].Properties.EM, n.Properties.Computed["width"])
			// 	cwV := utils.Max((colWidth-(n.Children[i].Properties.Computed["width"]+(b*2)))/2, 0)
			// 	xOffset = cwV
			// }

			posTrack := map[int]int{}
			p := 0

			for a, column := range orderedNode {
				var maxColumnHeight float32
				for _, item := range column {
					itemState := s[item.Properties.Id]
					maxColumnHeight = utils.Max(itemState.Height, maxColumnHeight)
				}

				yOffset = s[n.Children[0].Properties.Id].Y
				for _, item := range column {
					var i int
					for c, v := range n.Children {
						if v.Properties.Id == item.Properties.Id {
							i = c
						}
					}
					posTrack[p] = i
					p++
					itemState := s[item.Properties.Id]
					cState := s[n.Children[i].Properties.Id]
					// n.Children[i] = item
					if n.Style["justify-content"] == "space-between" {
						cwV := utils.Max((colWidth - (itemState.Width)), 0)
						// fmt.Println(colWidth, (itemState.Width), cwV, xOffset)
						if a == 0 {
							cState.X += 0
						} else if a == len(orderedNode)-1 {
							cState.X += cwV
						} else {
							cState.X += cwV / 2
						}
					} else if n.Style["justify-content"] == "flex-end" || n.Style["justify-content"] == "center" {
						dif := self.Width - (xOffset)
						if n.Style["justify-content"] == "center" {
							dif = dif / 2
						}
						cState.X += dif
					} else if n.Style["justify-content"] == "flex-start" || n.Style["justify-content"] == "" {
						cState.X += xOffset
					} else {
						cwV := utils.Max((colWidth-(itemState.Width))/2, 0)
						var offset float32
						if n.Style["justify-content"] == "space-evenly" {
							offset = ((cwV * 2) / float32(len(orderedNode))) * float32(a)
						}
						cState.X += xOffset + (cwV - offset)
					}
					cState.Y = yOffset
					yOffset += maxColumnHeight
					(*state)[n.Children[i].Properties.Id] = cState
					i++
				}
				xOffset += colWidth
			}

			content := n.Style["align-content"]

			if n.Style["flex-direction"] == "column" {
				content = n.Style["justify-content"]
			}

			if content != "" && content != "flex-start" {
				var min, max, rows, col, currY float32
				min = 1000000000000
				for _, v := range n.Children {
					vState := s[v.Properties.Id]
					min = utils.Min(min, vState.Y)
					max = utils.Max(max, vState.Height+vState.Y)
					if vState.Y > currY {
						rows++
						currY = vState.Y
					}
				}

				height := max - min
				rowHeight := ((self.Height - height) / rows)
				for e := range n.Children {
					i := posTrack[e]
					cState := s[n.Children[i].Properties.Id]
					row := float32(int(e % int(rows)))
					if row == 0 {
						col++
					}
					if len(orderedNode[int(col)-1]) <= int(row) {
						row = 0
					}

					if content == "center" {
						cState.Y += (self.Height - height) / 2
					} else if content == "flex-end" {
						cState.Y += (self.Height - height)
					} else if content == "space-around" {
						cState.Y += (rowHeight * row) + (rowHeight / 2)
					} else if content == "space-evenly" {
						cState.Y += (rowHeight * row) + (rowHeight / 2)
					} else if content == "space-between" {
						cState.Y += (((self.Height - height) / (rows - 1)) * row)
					} else if content == "stretch" {
						cState.Y += (rowHeight * row)
						if n.Children[i].Style["height"] == "" {
							cState.Height = self.Height / rows
						}
					}
					(*state)[n.Children[i].Properties.Id] = cState

				}

			}
			(*state)[n.Properties.Id] = self
		},
	}
}

func order(p element.Node, state *map[string]element.State, elements []element.Node, direction string, reversed, wrap bool) [][]element.Node {
	s := *state
	self := s[p.Properties.Id]
	var dir, marginStart, marginEnd string
	if direction == "column" {
		dir = "Height"
		marginStart = "Top"
		marginEnd = "Bottom"
	} else {
		dir = "Width"
		marginStart = "Left"
		marginEnd = "Right"
	}
	max, _ := utils.GetStructField(&self, dir)

	nodes := [][]element.Node{}

	if wrap {
		counter := 0
		if direction == "column" {
			collector := []element.Node{}
			for _, v := range elements {
				vState := s[v.Properties.Id]
				elMax := vState.Height
				elMS, _ := utils.GetStructField(&vState.Margin, marginStart)
				elME, _ := utils.GetStructField(&vState.Margin, marginEnd)
				tMax := elMax + elMS.(float32) + elME.(float32)
				if counter+int(tMax) < int(max.(float32)) {
					collector = append(collector, v)
				} else {
					if reversed {
						slices.Reverse(collector)
					}
					nodes = append(nodes, collector)
					collector = []element.Node{}
					collector = append(collector, v)
					counter = 0
				}
				counter += int(tMax)
			}
			if len(collector) > 0 {
				nodes = append(nodes, collector)
			}
		} else {
			var mod int
			for _, v := range elements {
				vState := s[v.Properties.Id]
				elMax := vState.Width
				elMS, _ := utils.GetStructField(&vState.Margin, marginStart)
				elME, _ := utils.GetStructField(&vState.Margin, marginEnd)
				tMax := elMax + elMS.(float32) + elME.(float32)
				if counter+int(tMax) < int(max.(float32)) {
					if len(nodes)-1 < mod {
						nodes = append(nodes, []element.Node{v})
					} else {
						nodes[mod] = append(nodes[mod], v)
					}
				} else {
					mod = 0
					counter = 0
					if len(nodes)-1 < mod {
						nodes = append(nodes, []element.Node{v})
					} else {
						nodes[mod] = append(nodes[mod], v)
					}
				}
				counter += int(tMax)
				mod++
			}
			if reversed {
				slices.Reverse(nodes)
			}
		}
	} else {
		var tMax float32
		for _, v := range elements {
			vState := s[v.Properties.Id]
			elMax, _ := utils.GetStructField(&vState, dir)
			elMS, _ := utils.GetStructField(&vState.Margin, marginStart)
			elME, _ := utils.GetStructField(&vState.Margin, marginEnd)
			tMax += elMax.(float32) + elMS.(float32) + elME.(float32)
		}

		pMax, _ := utils.GetStructField(&self, dir)

		// Resize node to fit
		var newSize float32
		if tMax > pMax.(float32) {
			newSize = pMax.(float32) / float32(len(elements))
		}
		if dir == "Width" {
			for _, v := range elements {
				vState := s[v.Properties.Id]
				if newSize != 0 {
					vState.Width = newSize
				}
				nodes = append(nodes, []element.Node{v})
			}
			if reversed {
				slices.Reverse(nodes)
			}
		} else {
			nodes = append(nodes, []element.Node{})
			for _, v := range elements {
				vState := s[v.Properties.Id]
				if newSize != 0 {
					vState.Height = newSize
				}
				nodes[0] = append(nodes[0], v)
			}
			if reversed {
				slices.Reverse(nodes[0])
			}
		}

	}
	(*state)[p.Properties.Id] = self

	return nodes
}
