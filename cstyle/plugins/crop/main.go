package crop

import (
	"gui/cstyle"
	"gui/element"
	"gui/library"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Selector: func(n *element.Node) bool {
			if n.Style["overflow"] != "" || n.Style["overflow-x"] != "" || n.Style["overflow-y"] != "" {
				// fmt.Println(n.Properties.Id)
				return true
			} else {
				return false
			}
		},
		Level: 1,
		Handler: func(n *element.Node, state *map[string]element.State, shelf *library.Shelf) {
			// !TODO: Needs to find crop bounds for X
			s := *state
			self := s[n.Properties.Id]
			// fmt.Println(n.Properties.Id)
			// if n.ScrollY == 0 {
			// 	return
			// }
			scrollY := findScroll(n)
			// !TODO: Limit scroll, and make scroll bar
			// + in the styles per what ever the css way to do it is

			minY, maxY := findBounds(n, state)
			// !ISSUE: Can add to the scroll value while it is pegged out
			scrollAmt := ((maxY - minY) + self.Padding.Bottom) / self.Height
			if scrollAmt > 1 {
				diff := scrollAmt - 1
				if scrollY > int(self.Height*diff) {
					scrollY = int(self.Height * diff)
					n.ScrollY = scrollY
				}
			}

			// fmt.Println(n.Properties.Id, n.PseudoElements)

			// DOnt work
			// !TODO: The width of the scroll bar needs to effect the width of the content inside of the container, needs
			// + to be moved up the render chain
			// + Also ::before and ::after can be handled by a plugin but will also need to be moved up

			// need to run before everything else, parents styles that need to change are width -= scrollbar width, padding-right += scrollbar-width

			// !TODO: Add keygen
			// key := n.Properties.Id + "123123"
			// exists := shelf.Check(key)
			// if exists {
			// lookup := make(map[string]struct{}, len(self.Textures))
			// for _, v := range self.Textures {
			// 	lookup[v] = struct{}{}
			// }

			// if _, found := lookup[key]; !found {
			// 	self.Textures = append(self.Textures, key)
			// }
			// } else {
			// width := int(self.Width + self.Padding.Left + self.Padding.Right)
			// height := int(self.Height + self.Padding.Top + self.Padding.Bottom)
			// ctx := canvas.NewCanvas(width, int(height))
			// ctx.FillStyle = color.RGBA{0, 255, 0, 255}
			// ctx.LineWidth = 1
			// ctx.BeginPath()
			// ctx.Rect(width-32, -1, 30, int(height)-1)
			// ctx.Fill()
			// ctx.ClosePath()
			// fmt.Println(self.Textures, ctx.Context.RGBAAt(width+10, 10))
			// self.Textures = append(self.Textures, shelf.Set(key, ctx.Context))
			// }
			// self.Width -= 30

			for _, v := range n.Children {
				if v.Style["position"] == "fixed" || v.TagName == "grim-scrollbar" {
					continue
				}
				child := s[v.Properties.Id]

				if (child.Y+child.Height)-float32(scrollY) < (self.Y) || (child.Y-float32(scrollY)) > self.Y+self.Height {
					child.Hidden = true
					(*state)[v.Properties.Id] = child
				} else {
					child.Hidden = false
					yCrop := 0
					height := int(child.Height)
					// !ISSUE: Text got messed up after the cropping? also in the raylib adapter with add the drawrect crop thing
					if child.Y-float32(scrollY) < (self.Y) {
						yCrop = int((self.Y) - (child.Y - float32(scrollY)))
						height = int(child.Height) - yCrop
					} else if (child.Y-float32(scrollY))+child.Height > self.Y+self.Height {
						diff := ((child.Y - float32(scrollY)) + child.Height) - (self.Y + self.Height)
						height = int(child.Height) - int(diff)
					}
					child.Crop = element.Crop{
						X:      0,
						Y:      yCrop,
						Width:  int(child.Width),
						Height: height,
					}
					(*state)[v.Properties.Id] = child
					// child.Y -= float32(n.ScrollY)

					updateChildren(v, state, scrollY)
				}
			}
			(*state)[n.Properties.Id] = self
		},
	}
}

func updateChildren(n *element.Node, state *map[string]element.State, offset int) {
	self := (*state)[n.Properties.Id]
	self.Y -= float32(offset)
	(*state)[n.Properties.Id] = self
	for _, v := range n.Children {
		updateChildren(v, state, offset)
	}
}

func findScroll(n *element.Node) int {
	if n.ScrollY != 0 {
		return n.ScrollY
	} else {
		for _, v := range n.Children {
			if v.Style["overflow"] == "" && v.Style["overflow-x"] == "" && v.Style["overflow-y"] == "" {
				s := findScroll(v)
				if s != 0 {
					return s
				}
			}
		}
		return 0
	}
}

func findBounds(n *element.Node, state *map[string]element.State) (float32, float32) {
	s := *state
	var minY, maxY float32
	minY = 10e10
	for _, v := range n.Children {
		child := s[v.Properties.Id]
		if child.Y < minY {
			minY = child.Y
		}
		if child.Y+child.Height > maxY {
			maxY = child.Y + child.Height
		}
		nMinY, nMaxY := findBounds(v, state)

		if nMinY < minY {
			minY = nMinY
		}
		if nMaxY > maxY {
			maxY = nMaxY
		}
	}
	return minY, maxY
}
