package flex

import (
	"fmt"
	"gui/cstyle"
	"gui/element"
	"gui/utils"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"display": "flex",
		},
		Level: 1,
		Handler: func(n *element.Node) {

			n.Width, _ = utils.ConvertToPixels("100%", n.EM, n.Parent.Width)

			fmt.Println(n.Id)
			siblings := n.Children

			fmt.Println(len(siblings))

			tW := float32(0)
			for _, v := range siblings {
				tW += v.Width
			}
			fmt.Println(n.Width, tW)
			for i, _ := range siblings {
				n.Children[i].X += ((n.Width - tW) / float32(len(siblings))) * float32(i)
			}
		},
	}
}
