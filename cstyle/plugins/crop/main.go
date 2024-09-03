package crop

import (
	"gui/cstyle"
	"gui/element"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Selector: func(n *element.Node) bool {
			if n.Style["overflow"] != "" || n.Style["overflow-x"] != "" || n.Style["overflow-y"] != "" {
				// fmt.Println(n.Properties.Id)
				return true
			} else {
				return false
			}
		},
		Level: 3,
		Handler: func(n *element.Node, state *map[string]element.State) {
			// !TODO: Needs to find crop bounds
			// + If you move the mouse it auto scrolls
			// + scrolling a child will cause the parent to scroll
			// + if an element is not overflowing no scrolling
			s := *state
			self := s[n.Properties.Id]
			// fmt.Println(n.Properties.Id)
			// if n.ScrollY == 0 {
			// 	return
			// }
			for _, v := range n.Children {
				if v.Style["position"] == "fixed" {
					continue
				}
				child := s[v.Properties.Id]
				if (child.Y+child.Height)-float32(n.ScrollY) < self.Y || (child.Y-float32(n.ScrollY)) > self.Y+self.Height {
					child.Hidden = true
					(*state)[v.Properties.Id] = child
				} else {
					child.Hidden = false
					(*state)[v.Properties.Id] = child
					// child.Y -= float32(n.ScrollY)
					updateChildren(v, state, n.ScrollY)
				}
			}

			(*state)[n.Properties.Id] = self
		},
	}
}

func updateChildren(n *element.Node, state *map[string]element.State, offset int) {
	self := (*state)[n.Properties.Id]
	self.Y -= float32(offset)
	(*state)[n.Properties.Id] = self
	for _, v := range n.Children {
		updateChildren(v, state, offset)
	}
}
