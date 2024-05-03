package textAlign

import (
	"fmt"
	"gui/cstyle"
	"gui/element"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"text-align": "*",
		},
		Level: 2,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]
			minX := float32(9e15)
			maxXW := float32(0)

			// fmt.Println(n.Properties.Id, len(n.Children))
			if self.Style["text-align"] == "center" {
				if len(n.Children) > 0 {
					baseX := s[n.Children[0].Properties.Id].X
					baseY := s[n.Children[0].Properties.Id].Y + s[n.Children[0].Properties.Id].Height
					last := 0
					for i := 1; i < len(n.Children)-1; i++ {

						cState := s[n.Children[i].Properties.Id]
						next := s[n.Children[i+1].Properties.Id]
						if baseY != next.Y+next.Height {
							baseY = next.Y + next.Height
							fmt.Println(last, i)
							for a := last; a < i; a++ {
								cState := s[n.Children[a].Properties.Id]
								cState.X += ((self.Width - (maxXW - minX)) / 2) - (baseX - self.X)
								(*state)[n.Children[a].Properties.Id] = cState
							}
							minX = 0
							maxXW = 9e15
							last = i
						}
						if cState.X < minX {
							minX = cState.X
						}
						if (cState.X + cState.Width) > maxXW {
							maxXW = cState.X + cState.Width
						}

					}

				}
			}

			(*state)[n.Properties.Id] = self
		},
	}
}
