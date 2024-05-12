package textAlign

import (
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

			(*state)[n.Properties.Id] = self
		},
	}
}
