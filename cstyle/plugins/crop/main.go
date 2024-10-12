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
			// parent := s[n.Parent.Properties.Id]

			scrollTop := findScroll(n)
			fmt.Println(scrollTop)
			minY, maxY := findBounds(n, state)

			// // !ISSUE: Can add to the scroll value while it is pegged out
			scrollAmt := ((maxY - minY) + self.Padding.Bottom + self.Padding.Top) / self.Height

			for _, v := range n.Children {
				if v.TagName == "grim-scrollbar" {
					if scrollAmt > 1 {
						diff := 1 - (scrollAmt - 1)
						p := s[v.Children[0].Properties.Id]
						p.Height = self.Height * diff
						p.Y += float32(scrollAmt)
						if self.Y+self.Height < p.Y+p.Height {
							p.Y = (self.Y + self.Height) - p.Height
						} else if p.Y < self.Y {
							p.Y = self.Y
						}

						(*state)[v.Children[0].Properties.Id] = p
					} else {
						p := s[v.Properties.Id]
						p.Hidden = true
						(*state)[v.Properties.Id] = p
						p = s[v.Children[0].Properties.Id]
						p.Hidden = true
						(*state)[v.Children[0].Properties.Id] = p
						// originalWidth, originalPaddingRight := getOriginalStyles(n.Style)
						// fmt.Println(n.Style["width"], originalWidth, originalPaddingRight)
						// if originalWidth != "" {
						// 	self.Width = utils.ConvertToPixels(originalWidth, self.EM, parent.Width)
						// } else {
						// 	self.Width = 0
						// }
						// if originalPaddingRight != "" {
						// 	self.Padding.Right = utils.ConvertToPixels(originalPaddingRight, self.EM, parent.Width)
						// } else {
						// 	self.Padding.Right = 0
						// }
						// fmt.Println(self.Width, self.Padding.Right)
					}
					break
				}
			}
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
	if n.ScrollTop != 0 {
		return n.ScrollTop
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

// // Helper function to extract the original value from a calc expression
// func extractOriginalValue(value string, operator string) string {
// 	if strings.HasPrefix(value, "calc(") && strings.Contains(value, operator) {
// 		// Trim "calc(" and everything after the operator
// 		start := strings.TrimPrefix(value, "calc(")
// 		end := strings.Index(start, operator)
// 		if end != -1 {
// 			return start[:end]
// 		}
// 	}
// 	return ""
// }

// // Function to return the original width and padding-right
// func getOriginalStyles(n map[string]string) (originalWidth, originalPaddingRight string) {
// 	// Get the current width and padding-right
// 	currentWidth := n["width"]
// 	pr := n["padding-right"]

// 	// Handle the case where padding-right might be in padding shorthand
// 	if pr == "" {
// 		if n["padding"] != "" {
// 			pr = n["padding"]
// 		}
// 	}

// 	// Extract the original width by looking for calc and -
// 	originalWidth = extractOriginalValue(currentWidth, "-")

// 	// Extract the original padding-right by looking for calc and +
// 	originalPaddingRight = extractOriginalValue(pr, "+")

// 	return originalWidth, originalPaddingRight
// }
