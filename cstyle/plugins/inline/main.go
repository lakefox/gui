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
			copyOfX := self.X
			copyOfY := self.Y
			xCollect := copyOfX
			// !ISSUE: Look at the if statements to see if they are properly selecting the correct elements
			for i, v := range n.Parent.Children {
				vState := s[v.Properties.Id]
				// if the element is the element being calculated currently
				if v.Style["display"] != "inline" {
					xCollect = 0
				} else {
					if v.Properties.Id == n.Properties.Id {
						// then if the x coordinate of the current element plus the x shift from its prevous siblings and the width minus 2 pixels
						// is greater than the width of the parent plus it's x value
						// and it is not the first element
						if self.X+xCollect+self.Width > parent.Width+parent.X && i > 0 {
							// We need to shift the element
							// then find the prevous sibling and add its height to the Y value to shift it downwards
							// shift the element to the base x (which should be the parent x value) and reset the xCollect
							fmt.Println(n.Properties.Id, "broke", xCollect, self.X, self.Width, parent.Width, parent.X)
							sibling := s[n.Parent.Children[i-1].Properties.Id]
							self.Y += sibling.Height
							self.X = copyOfX
							xCollect = copyOfX
						} else if i > 0 {
							fmt.Println(n.Properties.Id, n.InnerText, "did not break", xCollect, self.X, self.Width, parent.Width)
							self.X += xCollect
							fmt.Println(self.X)
						}
						break
					} else if vState.X+vState.Width > parent.Width+copyOfX && i > 0 {
						xCollect = copyOfX
					} else {
						xCollect += vState.Width + vState.X
						if vState.X+vState.Width+xCollect > parent.Width+parent.X {
							// If the elements are on the same line add the width of the elemenet
							if i > 0 {
								sibling := s[n.Parent.Children[i-1].Properties.Id]
								// !ISSUE: Commenting this out my help but can't decide
								if n.Parent.Children[i-1].Style["display"] != "inline" {
									self.Y += sibling.Height
								}
							}
							// !ISSUE: when not added and just set equal it shift one that collides but when added it breaks the line everytime
							if colliderDetection(vState, self) {
								xCollect += vState.Width
							} else {
								xCollect = 0
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
