package crop

import (
	"fmt"
	"gui/cstyle"
	"gui/element"
	"strings"
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
		Handler: func(n *element.Node, state *map[string]element.State) {
			// !TODO: Needs to find crop bounds for X
			s := *state
			self := s[n.Properties.Id]
			// fmt.Println(n.Properties.Id)
			// if n.ScrollTop == 0 {
			// 	return
			// }
			scrollY := n.ScrollY
			// scrollY := findScroll(n)
			// thumb := scrollEl.QuerySelector("grim-thumb")
			// fmt.Println("this", s[scrollEl.Properties.Id].Y, s[thumb.Properties.Id].Y)
			// scrollTop := int((s[thumb.Properties.Id].Y - s[scrollEl.Properties.Id].Y) + float32(scrollY))
			// // !TODO: Limit scroll
			// // + in the styles per what ever the css way to do it is
			// // fmt.Println(scrollTop, scrollEl.ScrollY)
			// if scrollTop < 0 {
			// 	scrollTop = 0
			// }
			minY, maxY := findBounds(n, state)
			// // !ISSUE: Can add to the scroll value while it is pegged out
			scrollAmt := ((maxY - minY) + self.Padding.Bottom + self.Padding.Top) / self.Height
			// // if there is more than 100% of parent
			// if scrollAmt > 1 {
			// 	diff := scrollAmt - 1
			// 	if scrollTop > int(self.Height*diff) {
			// 		scrollTop = int(self.Height * diff)
			// 		// need to stop the scroll value by dispatching an event to set the scroll
			// 	}
			// }

			// fmt.Println(n.Properties.Id, n.PseudoElements)
			scrollTop := 0
			for _, v := range n.Children {
				if v.TagName == "grim-scrollbar" {
					if scrollAmt > 1 {
						diff := 1 - (scrollAmt - 1)
						p := s[v.Children[0].Properties.Id]
						p.Height = self.Height * diff
						p.Y += float32(scrollY)
						if self.Y+self.Height < p.Y+p.Height {
							p.Y = (self.Y + self.Height) - p.Height
						} else if p.Y < self.Y {
							p.Y = self.Y
						}
						fmt.Println(p.Y)

						scrollTop = int((p.Y - self.Y))
						(*state)[v.Children[0].Properties.Id] = p
					} else {
						p := s[v.Properties.Id]
						p.Hidden = true
						(*state)[v.Properties.Id] = p
					}
					break
				}
			}
			fmt.Println(scrollTop, scrollY)
			for _, v := range n.Children {
				if v.Style["position"] == "fixed" || v.TagName == "grim-scrollbar" {
					continue
				}
				child := s[v.Properties.Id]

				if (child.Y+child.Height)-float32(scrollTop) < (self.Y) || (child.Y-float32(scrollTop)) > self.Y+self.Height {
					child.Hidden = true
					(*state)[v.Properties.Id] = child
				} else {
					child.Hidden = false
					yCrop := 0
					height := int(child.Height)
					// !ISSUE: Text got messed up after the cropping? also in the raylib adapter with add the drawrect crop thing
					if child.Y-float32(scrollTop) < (self.Y) {
						yCrop = int((self.Y) - (child.Y - float32(scrollTop)))
						height = int(child.Height) - yCrop
					} else if (child.Y-float32(scrollTop))+child.Height > self.Y+self.Height {
						diff := ((child.Y - float32(scrollTop)) + child.Height) - (self.Y + self.Height)
						height = int(child.Height) - int(diff)
					}
					child.Crop = element.Crop{
						X:      0,
						Y:      yCrop,
						Width:  int(child.Width),
						Height: height,
					}
					(*state)[v.Properties.Id] = child
					// child.Y -= float32(n.ScrollTop)

					updateChildren(v, state, scrollTop)
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

func findBounds(n *element.Node, state *map[string]element.State) (float32, float32) {
	s := *state
	var minY, maxY float32
	minY = 10e10
	for _, v := range n.Children {
		if strings.HasPrefix(v.TagName, "grim") {
			continue
		}
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

func findScroll(n *element.Node) int {
	// THis function should only return a value if none of its children have a value
	for _, v := range n.Children {
		if v.Style["overflow"] == "" && v.Style["overflow-x"] == "" && v.Style["overflow-y"] == "" {
			s := findScroll(v)
			if s != 0 {
				return 0
			}
		}
	}
	return n.ScrollY
}
