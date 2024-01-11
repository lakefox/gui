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

			n.Width, _ = utils.ConvertToPixels("100%", n.EM, n.Parent.Width)

			// Brief: justify does not align the bottom row correctly
			//        y axis also needs to be done
			offset(n, n.Styles["justify-content"], "X", "Width")
			// offset(n, n.Styles["align-content"], "Y", "Height")
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
					nodes[mod] = append(nodes[mod], v)
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
		if tMax > pMax.(float32) {
			newSize := pMax.(float32) / float32(len(elements))

			if dir == "Width" {
				for _, v := range elements {
					v.Width = newSize
					nodes = append(nodes, []element.Node{v})
				}
				if reversed {
					slices.Reverse(nodes)
				}
			} else {
				nodes = append(nodes, []element.Node{})
				for _, v := range elements {
					v.Height = newSize
					nodes[0] = append(nodes[0], v)
				}
				if reversed {
					slices.Reverse(nodes[0])
				}
			}
		}
	}

	return nodes
}

func offset(n *element.Node, styleKey, axis, side string) {
	nAxisInt, _ := utils.GetStructField(n, axis)
	var nAxis float32
	if nAxisInt != nil {
		nAxis = nAxisInt.(float32)
	}
	nSideInt, _ := utils.GetStructField(n, side)
	var nSide float32
	if nSideInt != nil {
		nSide = nSideInt.(float32)
	}
	var xOffset float32
	if styleKey == "" || styleKey == "flex-start" || styleKey == "flex-end" || styleKey == "center" {
		chunkStart := 0
		xOffset = nAxis
		for i, v := range n.Children {
			vSideInt, _ := utils.GetStructField(v, side)
			var vSide float32
			if vSideInt != nil {
				vSide = vSideInt.(float32)
			}
			vAxisInt, _ := utils.GetStructField(v, axis)
			var vAxis float32
			if vAxisInt != nil {
				vAxis = vAxisInt.(float32)
			}
			if xOffset+vSide > nSide {
				if styleKey == "flex-end" || styleKey == "center" {
					dif := nSide - (xOffset - nAxis)
					if styleKey == "center" {
						dif = dif / 2
					}
					for a := chunkStart; a < i; a++ {
						// n.Children[a].X += dif
						utils.SetStructFieldValue(n.Children[a], axis, vAxis+dif)
					}
				}
				chunkStart = i
				xOffset = nAxis
			}
			utils.SetStructFieldValue(n, axis, xOffset)
			xOffset += v.Width
		}
		if styleKey == "flex-end" || styleKey == "center" {
			dif := nSide - (xOffset)
			if styleKey == "center" {
				dif = dif / 2
			}
			for a := chunkStart; a < len(n.Children); a++ {
				vAxisInt, _ := utils.GetStructField(n.Children[a], axis)
				var vAxis float32
				if vAxisInt != nil {
					vAxis = vAxisInt.(float32)
				}
				utils.SetStructFieldValue(n.Children[a], axis, vAxis+dif)

			}
		}
	} else {
		verbs := strings.Split(n.Styles["flex-direction"], "-")
		orderedNode := order(*n, n.Children, verbs[0], len(verbs) > 1, n.Styles["flex-wrap"] == "wrap")

		var i int

		colWidth := n.Width / float32(len(orderedNode))

		if styleKey == "space-evenly" {
			cwV := utils.Max((colWidth-(n.Children[i].Width))/2, 0)
			xOffset = cwV
		}

		for a, column := range orderedNode {
			var maxColumnHeight float32
			for _, item := range column {
				maxColumnHeight = utils.Max(item.Height, maxColumnHeight)
			}

			for _, item := range column {
				n.Children[i] = item
				if styleKey == "space-between" {
					cwV := utils.Max((colWidth - (item.Width)), 0)
					if a == 0 {
						n.Children[i].X += xOffset
					} else if a == len(orderedNode)-1 {
						n.Children[i].X += xOffset + cwV
					} else {
						n.Children[i].X += xOffset + cwV/2
					}
				} else {
					cwV := utils.Max((colWidth-(item.Width))/2, 0)
					var offset float32
					if styleKey == "space-evenly" {
						offset = ((cwV * 2) / float32(len(orderedNode))) * float32(a)
					}
					n.Children[i].X += xOffset + (cwV - offset)

				}
				i++
			}
			xOffset += colWidth
		}
	}
}
