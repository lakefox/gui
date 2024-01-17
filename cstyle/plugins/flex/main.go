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
		Styles: map[string]string{
			"display":         "flex",
			"justify-content": "*",
			"align-content":   "*",
			"flex-wrap":       "*",
			"flex-direction":  "*",
		},
		Level: 1,
		Handler: func(n *element.Node) {

			// Brief: justify does not align the bottom row correctly
			//        y axis also needs to be done
			verbs := strings.Split(n.Style["flex-direction"], "-")
			orderedNode := order(*n, n.Children, verbs[0], len(verbs) > 1, n.Style["flex-wrap"] == "wrap")

			var i int

			colWidth := n.Width / float32(len(orderedNode))

			var xOffset, yOffset float32
			if n.Style["justify-content"] == "space-evenly" {
				b, _ := utils.ConvertToPixels(n.Children[i].Border.Width, n.Children[i].EM, n.Width)
				cwV := utils.Max((colWidth-(n.Children[i].Width+(b*2)))/2, 0)
				xOffset = cwV
			}

			for a, column := range orderedNode {
				var maxColumnHeight float32
				for _, item := range column {
					maxColumnHeight = utils.Max(item.Height, maxColumnHeight)
				}

				yOffset = n.Children[0].Y
				for _, item := range column {
					n.Children[i] = item
					if n.Style["justify-content"] == "space-between" {
						cwV := utils.Max((colWidth - (item.Width)), 0)
						if a == 0 {
							n.Children[i].X += xOffset
						} else if a == len(orderedNode)-1 {
							n.Children[i].X += xOffset + cwV
						} else {
							n.Children[i].X += xOffset + cwV/2
						}
					} else if n.Style["justify-content"] == "flex-end" || n.Style["justify-content"] == "center" {
						dif := n.Width - (xOffset)
						if n.Style["justify-content"] == "center" {
							dif = dif / 2
						}
						n.Children[i].X += dif
					} else if n.Style["justify-content"] == "flex-start" || n.Style["justify-content"] == "" {
						n.Children[i].X += xOffset
					} else {
						cwV := utils.Max((colWidth-(item.Width))/2, 0)
						var offset float32
						if n.Style["justify-content"] == "space-evenly" {
							offset = ((cwV * 2) / float32(len(orderedNode))) * float32(a)
						}
						n.Children[i].X += xOffset + (cwV - offset)

					}
					n.Children[i].Y = yOffset
					yOffset += maxColumnHeight
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
					min = utils.Min(min, v.Y)
					max = utils.Max(max, v.Height+v.Y)
					if v.Y > currY {
						rows++
						currY = v.Y
					}
				}

				height := max - min
				rowHeight := ((n.Height - height) / rows)

				for i := range n.Children {
					row := float32(int(i % int(rows)))
					if row == 0 {
						col++
					}
					if len(orderedNode[int(col)-1]) <= int(row) {
						row = 0
					}

					if content == "center" {
						n.Children[i].Y += (n.Height - height) / 2
					} else if content == "flex-end" {
						n.Children[i].Y += (n.Height - height)
					} else if content == "space-around" {
						n.Children[i].Y += (rowHeight * row) + (rowHeight / 2)
					} else if content == "space-evenly" {
						n.Children[i].Y += (rowHeight * row) + (rowHeight / 2)
					} else if content == "space-between" {
						n.Children[i].Y += (((n.Height - height) / (rows - 1)) * row)
					} else if content == "stretch" {
						n.Children[i].Y += (rowHeight * row)
						if n.Children[i].Style["height"] == "" {
							n.Children[i].Height = n.Height / rows
						}
					}
				}
			}

		},
	}
}

func order(p element.Node, elements []element.Node, direction string, reversed, wrap bool) [][]element.Node {
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
	max, _ := utils.GetStructField(&p, dir)

	nodes := [][]element.Node{}

	if wrap {
		counter := 0
		if direction == "column" {
			collector := []element.Node{}
			for _, v := range elements {
				elMax, _ := utils.GetStructField(&v, "Height")
				elMS, _ := utils.GetStructField(&v.Margin, marginStart)
				elME, _ := utils.GetStructField(&v.Margin, marginEnd)
				tMax := elMax.(float32) + elMS.(float32) + elME.(float32)
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
				elMax, _ := utils.GetStructField(&v, "Width")
				elMS, _ := utils.GetStructField(&v.Margin, marginStart)
				elME, _ := utils.GetStructField(&v.Margin, marginEnd)
				tMax := elMax.(float32) + elMS.(float32) + elME.(float32)
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
			elMax, _ := utils.GetStructField(&v, dir)
			elMS, _ := utils.GetStructField(&v.Margin, marginStart)
			elME, _ := utils.GetStructField(&v.Margin, marginEnd)
			tMax += elMax.(float32) + elMS.(float32) + elME.(float32)
		}

		pMax, _ := utils.GetStructField(&p, dir)

		// Resize node to fit
		var newSize float32
		if tMax > pMax.(float32) {
			newSize = pMax.(float32) / float32(len(elements))
		}
		if dir == "Width" {
			for _, v := range elements {
				if newSize != 0 {
					v.Width = newSize
				}
				nodes = append(nodes, []element.Node{v})
			}
			if reversed {
				slices.Reverse(nodes)
			}
		} else {
			nodes = append(nodes, []element.Node{})
			for _, v := range elements {
				if newSize != 0 {
					v.Height = newSize
				}
				nodes[0] = append(nodes[0], v)
			}
			if reversed {
				slices.Reverse(nodes[0])
			}
		}

	}

	return nodes
}
