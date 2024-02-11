package block

import (
	"gui/cstyle"
	"gui/element"
	"gui/utils"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"display": "block",
		},
		Level: 1,
		Handler: func(n *element.Node) {
			// If the element is display block and the width is unset then make it 100%

			if n.Style["width"] == "" {
				n.Properties.Width, _ = utils.ConvertToPixels("100%", n.Properties.EM, n.Parent.Properties.Width)
				m := utils.GetMP(*n, "margin")
				n.Properties.Width -= m.Right + m.Left
			} else {
				p := utils.GetMP(*n, "padding")
				n.Properties.Width += p.Right + p.Left
			}
		},
	}
}
