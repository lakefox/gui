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
			"display":         "flex",
			"justify-content": "*",
		},
		Level: 1,
		Handler: func(n *element.Node) {

			n.Width, _ = utils.ConvertToPixels("100%", n.EM, n.Parent.Width)

			fmt.Println(n.Id)
			siblings := n.Children

			fmt.Println(len(siblings))

			tW := float32(0)
			for _, v := range siblings {
				tW += v.Width + v.Margin.Right + v.Margin.Left
			}
			if tW > n.Width {
				fmt.Println("TOO BIG", tW/float32(len(siblings)))
				newSize := n.Width / float32(len(siblings))
				tW = n.Width
				for i := range siblings {
					n.Children[i].Width = newSize
				}
			}
			fmt.Println(n.Width, tW)

			fmt.Println(n.Styles["justify-content"])
			y := siblings[0].Y
			for i := range siblings {
				fmt.Println(n.Children[i].Width)
				n.Children[i].X += (((n.Width - tW) / float32(len(siblings))) * float32(i)) + (n.Children[i].Width * float32(i))
				n.Children[i].Y = y
			}
		},
	}
}
