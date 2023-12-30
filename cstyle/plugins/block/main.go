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
		Level: 0,
		Handler: func(n *element.Node) {
			// If the element is display block and the width is unset then make it 100%
			if n.Styles["width"] == "" {
				n.Width, _ = utils.ConvertToPixels("100%", n.EM, n.Parent.Width)
				n.Width -= n.Margin.Right + n.Margin.Left
			} else {
				n.Width += n.Padding.Right + n.Padding.Left
			}
		},
	}
}
