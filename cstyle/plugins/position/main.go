package position

import (
	"gui/cstyle"
	"gui/element"
	"gui/utils"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"position": "*",
		},
		Level: 0,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]
			parent := s[n.Parent.Properties.Id]

			width, height := self.Width, self.Height
			x, y := self.X, self.Y

			var top, left, right, bottom bool = false, false, false, false

			m := utils.GetMP(*n, "margin")

			if n.Style["position"] == "absolute" {
				bas := utils.GetPositionOffsetNode(n)
				base := s[bas.Properties.Id]
				if n.Style["top"] != "" {
					v, _ := utils.ConvertToPixels(n.Style["top"], self.EM, parent.Width)
					y = v + base.Y
					top = true
				}
				if n.Style["left"] != "" {
					v, _ := utils.ConvertToPixels(n.Style["left"], self.EM, parent.Width)
					x = v + base.X
					left = true
				}
				if n.Style["right"] != "" {
					v, _ := utils.ConvertToPixels(n.Style["right"], self.EM, parent.Width)
					x = (base.Width - width) - v
					right = true
				}
				if n.Style["bottom"] != "" {
					v, _ := utils.ConvertToPixels(n.Style["bottom"], self.EM, parent.Width)
					y = (base.Height - height) - v
					bottom = true
				}
			} else {
				for i, v := range n.Parent.Children {
					if v.Properties.Id == n.Properties.Id {
						if i-1 > 0 {
							sib := n.Parent.Children[i-1]
							sibling := s[sib.Properties.Id]
							if n.Style["display"] == "inline" {
								if sib.Style["display"] == "inline" {
									y = sibling.Y
								} else {
									y = sibling.Y + sibling.Height
								}
							} else {
								y = sibling.Y + sibling.Height
							}
						}
						break
					} else if n.Style["display"] != "inline" {
						mc := utils.GetMP(v, "margin")
						pc := utils.GetMP(v, "padding")
						vState := s[v.Properties.Id]
						y += mc.Top + mc.Bottom + pc.Top + pc.Bottom + vState.Height
					}
				}
			}

			// Display modes need to be calculated here

			relPos := !top && !left && !right && !bottom

			if left || relPos {
				x += m.Left
			}
			if top || relPos {
				y += m.Top
			}
			if right {
				x -= m.Right
			}
			if bottom {
				y -= m.Bottom
			}

			self.X = x
			self.Y = y
			self.Width = width
			self.Height = height
		},
	}
}
