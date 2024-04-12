package inline

import (
	"fmt"
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
			fmt.Println("start: ", n.InnerText)
			copyOfX := self.X
			baseY := s[n.Parent.Children[0].Properties.Id].Y
			xCollect := float32(0)
			for i, v := range n.Parent.Children {
				vState := s[v.Properties.Id]
				if v.Properties.Id == n.Properties.Id {
					fmt.Println(23)
					if self.X+self.Width-2 > parent.Width+copyOfX && i > 0 {
						fmt.Println(25)
						sibling := s[n.Parent.Children[i-1].Properties.Id]
						self.Y += sibling.Height
						self.X = copyOfX + self.Width
					}
					if i > 0 {
						fmt.Println(31)
						if n.Parent.Children[i-1].Style["display"] == "inline" {
							fmt.Println(33)
							sibling := s[n.Parent.Children[i-1].Properties.Id]

							if sibling.Text.X+self.Text.Width < int(sibling.Width) {
								fmt.Println(37)
								self.Y -= float32(sibling.Text.LineHeight)
								self.X += float32(sibling.Text.X)
							}
						}
					}
					break
				} else if v.Style["display"] == "inline" {
					fmt.Println(45)
					self.X += xCollect + vState.Width
					self.Y = baseY
				} else {
					fmt.Println(49)
					self.X = copyOfX
				}
			}
			(*state)[n.Properties.Id] = self
		},
	}
}
