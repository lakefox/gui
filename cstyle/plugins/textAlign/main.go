package textAlign

import (
	"gui/cstyle"
	"gui/element"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"text-align": "*",
			"inlineText": "*",
		},
		Level: 2,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]
			// parent := s[n.Parent.Properties.Id]
			minX := float32(9e15)
			maxXW := float32(0)

			// fmt.Println(n.Properties.Id)

			// fmt.Println(n.Properties.Id, len(n.Children))
			if self.Style["text-align"] == "center" {
				if len(n.Children) > 0 {
					minX = s[n.Children[0].Properties.Id].X
					baseY := s[n.Children[0].Properties.Id].Y + s[n.Children[0].Properties.Id].Height
					last := 0
					for i := 1; i < len(n.Children)-1; i++ {

						cState := s[n.Children[i].Properties.Id]
						next := s[n.Children[i+1].Properties.Id]

						if cState.X < minX {
							minX = cState.X
						}
						if (cState.X + cState.Width) > maxXW {
							maxXW = cState.X + cState.Width
						}

						if baseY != next.Y+next.Height {
							baseY = next.Y + next.Height
							for a := last; a < i+1; a++ {
								cState := s[n.Children[a].Properties.Id]
								cState.X += (((self.Width + (self.Border.Width * 2)) - (maxXW - minX)) / 2) - (minX - self.X)
								(*state)[n.Children[a].Properties.Id] = cState
							}
							minX = 9e15
							maxXW = 0
							last = i + 1
						}

					}
					minX = s[n.Children[last].Properties.Id].X
					maxXW = s[n.Children[len(n.Children)-1].Properties.Id].X + s[n.Children[len(n.Children)-1].Properties.Id].Width
					for a := last; a < len(n.Children); a++ {
						cState := s[n.Children[a].Properties.Id]
						cState.X += (((self.Width + (self.Border.Width * 2)) - (maxXW - minX)) / 2) - (minX - self.X)
						(*state)[n.Children[a].Properties.Id] = cState
					}
				}
			} else if self.Style["text-align"] == "right" {
				if len(n.Children) > 0 {
					minX = s[n.Children[0].Properties.Id].X
					baseY := s[n.Children[0].Properties.Id].Y + s[n.Children[0].Properties.Id].Height
					last := 0
					for i := 1; i < len(n.Children)-1; i++ {

						cState := s[n.Children[i].Properties.Id]
						next := s[n.Children[i+1].Properties.Id]

						if cState.X < minX {
							minX = cState.X
						}
						if (cState.X + cState.Width) > maxXW {
							maxXW = cState.X + cState.Width
						}

						if baseY != next.Y+next.Height {
							baseY = next.Y + next.Height
							for a := last; a < i+1; a++ {
								cState := s[n.Children[a].Properties.Id]
								cState.X += ((self.Width + (self.Border.Width * 2)) - (maxXW - minX)) + ((self.X - minX) * 2)
								(*state)[n.Children[a].Properties.Id] = cState
							}
							minX = 9e15
							maxXW = 0
							last = i + 1
						}

					}
					minX = s[n.Children[last].Properties.Id].X
					maxXW = s[n.Children[len(n.Children)-1].Properties.Id].X + s[n.Children[len(n.Children)-1].Properties.Id].Width
					for a := last; a < len(n.Children); a++ {
						cState := s[n.Children[a].Properties.Id]
						cState.X += ((self.Width + (self.Border.Width * 2)) - (maxXW - minX)) + ((self.X - minX) * 2)
						(*state)[n.Children[a].Properties.Id] = cState
					}

				}
			}

			// if self.Style["inlineText"] == "true" {
			// 	if n.Properties.Id == n.Parent.Children[len(n.Parent.Children)-1].Properties.Id {
			// 		baseY := s[n.Parent.Children[0].Properties.Id].Y
			// 		b := 0

			// 		for a := 0; a < len(n.Parent.Children)-1; a++ {
			// 			next := s[n.Parent.Children[a+1].Properties.Id]
			// 			if next.Y != baseY {
			// 				if len(n.Parent.Children[a].InnerText) > 0 {
			// 					// this is the end of the line
			// 					// fmt.Println(b, a)
			// 					// fmt.Println("#####")
			// 					tallest := float32(0)
			// 					for i := b; i < a+1; i++ {
			// 						vState := s[n.Parent.Children[i].Properties.Id]
			// 						if vState.Y+vState.Height > tallest {
			// 							tallest = vState.Y + vState.Height
			// 						}
			// 					}
			// 					// fmt.Println("tallest", tallest)
			// 					for i := b; i < a+1; i++ {
			// 						vState := s[n.Parent.Children[i].Properties.Id]
			// 						// fmt.Println(tallest-(vState.Y+vState.Height), n.Parent.Children[i].InnerText, vState.Y)
			// 						vState.Y += float32(i * 10)
			// 						// vState.Y += tallest - (vState.Y + vState.Height)
			// 						(*state)[n.Parent.Children[i].Properties.Id] = vState
			// 					}

			// 					// fmt.Print("\n")
			// 				}
			// 				b = a + 1
			// 				baseY = next.Y
			// 			}
			// 			// !ISSUE: Last line not found
			// 		}

			// 		// fmt.Println(n.InnerText, self.Y, len(n.Parent.Children))
			// 	}
			// }

			(*state)[n.Properties.Id] = self
		},
	}
}

// func alignText(n *element.Node, s map[string]element.State, i int, state *map[string]element.State) {
// 	tallest := float32(0)
// 	endex := 0
// 	for a := i; a > 0; a-- {
// 		if s[n.Parent.Children[a].Properties.Id].Y != s[n.Parent.Children[i-1].Properties.Id].Y {
// 			endex = a
// 			break
// 		} else {
// 			tallest = utils.Max(tallest, s[n.Parent.Children[a].Properties.Id].Height)
// 		}
// 	}
// 	// !ISSUE: Find a better way then -4
// 	for a := i; a > endex-1; a-- {
// 		p := (*state)[n.Parent.Children[a].Properties.Id]
// 		if p.Height != tallest {
// 			p.Y += (tallest - p.Height) - 5
// 			(*state)[n.Parent.Children[a].Properties.Id] = p
// 		}
// 	}
// 	for a := i; a < len(n.Parent.Children)-1; a++ {
// 		p := (*state)[n.Parent.Children[a].Properties.Id]
// 		p.Y += 7
// 		(*state)[n.Parent.Children[a].Properties.Id] = p
// 	}
// }
