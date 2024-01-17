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
		Level: 0,
		Handler: func(n *element.Node) {
			copyOfX := n.X
			for i, v := range n.Parent.Children {
				if v.Id == n.Id {
					if n.X+n.Width-2 > n.Parent.Width+copyOfX && i > 0 {
						n.Y += float32(n.Parent.Children[i-1].Height)
						n.X = copyOfX
					}
					if i > 0 {
						if n.Parent.Children[i-1].Style["display"] == "inline" {
							if n.Parent.Children[i-1].Text.X+n.Text.Width < int(n.Parent.Children[i-1].Width) {
								n.Y -= float32(n.Parent.Children[i-1].Text.LineHeight)
								n.X += float32(n.Parent.Children[i-1].Text.X)
							}
						}
					}
					break
				} else if v.Style["display"] == "inline" {
					n.X += v.Width
				} else {
					n.X = copyOfX
				}
			}
		},
	}
}
