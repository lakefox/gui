package flex

import (
	"fmt"
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
			"flex-wrap":       "*",
			"flex-direction":  "*",
		},
		Level: 1,
		Handler: func(n *element.Node) {

			n.Width, _ = utils.ConvertToPixels("100%", n.EM, n.Parent.Width)

			fmt.Println(n.Id)

			verbs := strings.Split(n.Styles["flex-direction"], "-")

			orderedNode := order(*n, n.Children, verbs[0], len(verbs) > 1, n.Styles["flex-wrap"] == "wrap")

			fmt.Println("SIASDJAOISJ")
			for i, v := range orderedNode {
				fmt.Println(i, len(v))
			}

			for _, row := range orderedNode {
				for _, column := range row {
					// add align-content justify-content and align-items logic here
					// the elements are positioned correctly in the ordered array just need to calculate the px positions
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
		collector := []element.Node{}
		for _, v := range elements {
			elMax, _ := utils.GetStructField(&v, dir)
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
