package flex

import (
	"fmt"
	"gui/cstyle"
	"gui/element"
	"gui/utils"
	"strings"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Selector: func(n *element.Node) bool {
			styles := map[string]string{
				"display": "flex",
			}
			matches := true
			for name, value := range styles {
				if n.Style[name] != value && !(value == "*") && n.Style[name] != "" {
					matches = false
				}
			}
			return matches
		},
		Level: 1,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]

			verbs := strings.Split(n.Style["flex-direction"], "-")
			flexDirection := verbs[0]
			flexReversed := false
			if len(verbs) > 1 {
				flexReversed = true
			}

			flexWrapped := !(n.Style["flex-wrap"] != "nowrap")

			hAlign := n.Style["align-content"]
			if hAlign == "" {
				hAlign = "normal"
			}
			vAlign := n.Style["align-items"]
			if vAlign == "" {
				vAlign = "normal"
			}
			justify := n.Style["justify-items"]
			if justify == "" {
				justify = "normal"
			}

			fmt.Println(flexDirection, flexReversed, flexWrapped, hAlign, vAlign, justify)

			if flexDirection == "row" && !flexReversed && !flexWrapped {
				for _, v := range n.Children {
					vState := s[v.Properties.Id]

					fmt.Println(getMinSize(&v, vState))
					fmt.Println(getInnerSize(&v, s))
					pw, ph := getMinSize(&v, vState)
					w, h := getInnerSize(&v, s)
					w += pw
					h += ph
					// !INSIGHT: if you change the width of the element then you need to re line the text...
					vState.Width = 100
					// vState.Width = w
					vState.Height = h
					(*state)[v.Properties.Id] = vState
				}

			}

			(*state)[n.Properties.Id] = self
		},
	}
}

func getMinSize(n *element.Node, s element.State) (float32, float32) {
	minW := s.Padding.Left + s.Padding.Right
	minH := s.Padding.Top + s.Padding.Bottom
	if n.Style["width"] != "" {
		minW = s.Width
	}
	if n.Style["height"] != "" {
		minH = s.Height
	}
	return minW, minH
}

func getInnerSize(n *element.Node, s map[string]element.State) (float32, float32) {
	minx := float32(10e10)
	maxw := float32(0)
	miny := float32(10e10)
	maxh := float32(0)
	for _, v := range n.Children {
		vState := s[v.Properties.Id]
		minx = utils.Min(vState.X, minx)
		miny = utils.Min(vState.Y, miny)

		maxw = utils.Max(vState.X+vState.Width, maxw)
		maxh = utils.Max(vState.Y+vState.Height, maxh)
	}
	return maxw - minx, maxh - miny
}
