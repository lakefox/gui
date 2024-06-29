package inline

import (
	"gui/cstyle"
	"gui/element"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Selector: func(n *element.Node) bool {
			styles := map[string]string{
				"display": "inline",
			}
			matches := true
			for name, value := range styles {
				if n.Style[name] != value && !(value == "*") && n.Style[name] != "" {
					matches = false
				}
			}
			return matches
		},
		Level: 1,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]
			parent := s[n.Parent.Properties.Id]
			copyOfX := self.X
			copyOfY := self.Y

			// xCollect := float32(0)
			for i, v := range n.Parent.Children {
				// vState := s[v.Properties.Id]
				if i > 0 {
					if v.Style["position"] != "absolute" {
						if v.Properties.Id == n.Properties.Id {
							sib := n.Parent.Children[i-1]
							sibling := s[sib.Properties.Id]
							if sibling.X+sibling.Width+self.Width > ((parent.Width)-parent.Padding.Right)+parent.X {
								// Break Node.Id
								self.Y = sibling.Y + sibling.Height
								self.X = copyOfX
							} else {
								// Node did not break
								if sib.Style["display"] != "inline" {
									self.Y = sibling.Y + sibling.Height
								} else {
									self.Y = sibling.Y
									self.X = sibling.X + sibling.Width
								}
							}
							break
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
