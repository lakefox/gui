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
		Handler: func(n element.Node) element.Node {
			// If the element is display block and the width is unset then make it 100%
			if n.Style["width"] == "" {
				n.Properties.Width, _ = utils.ConvertToPixels("100%", n.Properties.EM, n.Parent.Properties.Width)
				n.Properties.Width -= n.Properties.Margin.Right + n.Properties.Margin.Left
			} else {
				n.Properties.Width += n.Properties.Padding.Right + n.Properties.Padding.Left
			}
			return n
		},
	}
}
