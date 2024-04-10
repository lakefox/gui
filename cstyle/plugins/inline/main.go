package inline

import (
	"gui/cstyle"
	"gui/element"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"display": "inline",
		},
		Level: 1,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]
			parent := s[n.Parent.Properties.Id]

			copyOfX := self.X
			for i, v := range n.Parent.Children {
				if v.Properties.Id == n.Properties.Id {
					if self.X+self.Width-2 > parent.Width+copyOfX && i > 0 {
						sibling := s[n.Parent.Children[i-1].Properties.Id]
						self.Y += sibling.Height
						self.X = copyOfX
					}
					if i > 0 {
						if n.Parent.Children[i-1].Style["display"] == "inline" {
							sibling := s[n.Parent.Children[i-1].Properties.Id]

							if sibling.Text.X+self.Text.Width < int(sibling.Width) {
								self.Y -= float32(sibling.Text.LineHeight)
								self.X += float32(sibling.Text.X)
							}
						}
					}
					break
				} else if v.Style["display"] == "inline" {
					vState := s[v.Properties.Id]
					self.X += vState.Width
				} else {
					self.X = copyOfX
				}
			}
			(*state)[n.Properties.Id] = self
		},
	}
}
