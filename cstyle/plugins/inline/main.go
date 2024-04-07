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
		Handler: func(n *element.Node) {
			copyOfX := n.Properties.X
			for i, v := range n.Parent.Children {
				if v.Properties.Id == n.Properties.Id {
					if n.Properties.X+n.Properties.Computed["width"]-2 > n.Parent.Properties.Computed["width"]+copyOfX && i > 0 {
						n.Properties.Y += float32(n.Parent.Children[i-1].Properties.Computed["height"])
						n.Properties.X = copyOfX
					}
					if i > 0 {
						if n.Parent.Children[i-1].Style["display"] == "inline" {
							if n.Parent.Children[i-1].Properties.Text.X+n.Properties.Text.Width < int(n.Parent.Children[i-1].Properties.Computed["width"]) {
								n.Properties.Y -= float32(n.Parent.Children[i-1].Properties.Text.LineHeight)
								n.Properties.X += float32(n.Parent.Children[i-1].Properties.Text.X)
							}
						}
					}
					break
				} else if v.Style["display"] == "inline" {
					n.Properties.X += v.Properties.Computed["width"]
				} else {
					n.Properties.X = copyOfX
				}
			}
		},
	}
}
