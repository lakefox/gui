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
			xCollect := float32(0)
			for i, v := range n.Parent.Children {
				vState := s[v.Properties.Id]
				if v.Properties.Id == n.Properties.Id {
					if self.X+xCollect+copyOfX+self.Width-2 > parent.Width+copyOfX && i > 0 {
						sibling := s[n.Parent.Children[i-1].Properties.Id]
						self.Y += sibling.Height
						self.X = copyOfX
						xCollect = copyOfX
					} else if i > 0 {
						self.X += xCollect
					}
					break
				} else if vState.X+xCollect+vState.Width-2 > parent.Width+copyOfX && i > 0 {
					// !ISSUE: Added the + xCollect and things got better but don't know why
					xCollect = copyOfX + xCollect
				} else {
					xCollect += vState.Width
				}
			}
			(*state)[n.Properties.Id] = self
		},
	}
}
