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
			copyOfY := self.Y
			xCollect := float32(0)
			for i, v := range n.Parent.Children {
				vState := s[v.Properties.Id]
				if vState.Style["position"] != "absolute" {
					if vState.Style["display"] != "inline" {
						xCollect = 0
					} else {
						if v.Properties.Id == n.Properties.Id {
							if self.X+xCollect+self.Width > ((parent.Width)-parent.Padding.Right)+parent.X && i > 0 {
								// Break Node
								sibling := s[n.Parent.Children[i-1].Properties.Id]
								self.Y += sibling.Height
								self.X = copyOfX
							} else if i > 0 {
								// Node did not break
								self.X += xCollect
							}
							break
						} else {
							if vState.Style["display"] == "inline" {
								if colliderDetection(vState, self) {
									xCollect += vState.Width
								} else {
									xCollect = 0
								}
							}
						}
					}
				}
			}
			propagateOffsets(n, copyOfX, copyOfY, self, state)
			(*state)[n.Properties.Id] = self
		},
	}
}

func propagateOffsets(n *element.Node, copyOfX, copyOfY float32, self element.State, state *map[string]element.State) {
	s := *state
	for _, v := range n.Children {
		vState := s[v.Properties.Id]
		vState.X += self.X - copyOfX
		vState.Y += self.Y - copyOfY
		if len(v.Children) > 0 {
			propagateOffsets(&v, copyOfX, copyOfY, self, state)
		}
		(*state)[v.Properties.Id] = vState
	}
}

func colliderDetection(s1, s2 element.State) bool {
	s1Min := s1.Y
	s1Max := s1.Y + s1.Height
	s2Min := s2.Y
	s2Max := s2.Y + s2.Height
	return s1Min > s2Min && s1Min < s2Max || s1Max > s2Min && s1Min < s2Max || s2Min > s1Min && s2Min < s1Max || s2Max > s1Min && s2Min < s1Max
}
