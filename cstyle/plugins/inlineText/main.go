package inlineText

import (
	"gui/cstyle"
	"gui/element"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"inlineText": "*",
		},
		Level: 3,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]
			// parent := s[n.Parent.Properties.Id]
			// fmt.Println("#######")
			// if self.Style["inlineText"] == "true" {
			// 	baseY := s[n.Parent.Children[0].Properties.Id].Y + s[n.Parent.Children[0].Properties.Id].Height
			// 	index := -1
			// 	for a := 0; a < len(n.Parent.Children)-1; a++ {
			// 		if n.Properties.Id == n.Parent.Children[a].Properties.Id {
			// 			index = a
			// 			break
			// 		}
			// 	}
			// 	lastLB := 0
			// 	for a := 0; a < len(n.Parent.Children)-1; a++ {
			// 		next := s[n.Parent.Children[a+1].Properties.Id]
			// 		if next.Y+next.Height != baseY {
			// 			fmt.Println(n.Parent.Children[a+1].Properties.Id, n.Parent.Children[a+1].InnerText)
			// 			// We are waiting to find the line break after the current word here

			// 			// but we need to calculat the offsets for every other one
			// 			if a > index {
			// 				tallest := float32(0)
			// 				for b := lastLB; b < a; b++ {
			// 					vState := s[n.Parent.Children[b].Properties.Id]
			// 					if vState.Y+vState.Height > tallest {
			// 						tallest = vState.Y + vState.Height
			// 					}
			// 				}
			// 				self.Y += tallest - (self.Y + self.Height)
			// 				break
			// 			}
			// 			lastLB = a
			// 			baseY = next.Y + next.Height
			// 		}
			// 		// !ISSUE: Last line not found
			// 	}

			// if n.Properties.Id == n.Parent.Children[len(n.Parent.Children)-1].Properties.Id {

			// 	b := 0

			// 	fmt.Println(n.InnerText)

			// 	for a := 0; a < len(n.Parent.Children)-1; a++ {
			// 		next := s[n.Parent.Children[a+1].Properties.Id]
			// 		if next.Y != baseY {
			// 			// if len(n.Parent.Children[a].InnerText) > 0 {
			// 			// this is the end of the line
			// 			// fmt.Println(b, a)
			// 			fmt.Println("#####")
			// 			// tallest := float32(0)
			// 			for i := b; i < a+1; i++ {
			// 				vState := s[n.Parent.Children[i].Properties.Id]
			// 				vState.Y += 10 * float32(i)
			// 				(*state)[n.Parent.Children[i].Properties.Id] = vState
			// 				// if vState.Y+vState.Height > tallest {
			// 				// 	tallest = vState.Y + vState.Height
			// 				// }
			// 				fmt.Println(n.Parent.Children[i].InnerText)
			// 			}
			// 			// // fmt.Println("tallest", tallest)
			// 			// for i := b; i < a+1; i++ {
			// 			// 	vState := s[n.Parent.Children[i].Properties.Id]
			// 			// 	// fmt.Println(tallest-(vState.Y+vState.Height), n.Parent.Children[i].InnerText, vState.Y)
			// 			// 	vState.Y += float32(i * 10)
			// 			// 	// vState.Y += tallest - (vState.Y + vState.Height)
			// 			// 	(*state)[n.Parent.Children[i].Properties.Id] = vState
			// 			// }

			// 			// fmt.Print("\n")
			// 			// }
			// 			b = a + 1
			// 			baseY = next.Y
			// 		}
			// 		// !ISSUE: Last line not found
			// 	}

			// 	// fmt.Println(n.InnerText, self.Y, len(n.Parent.Children))
			// }
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
