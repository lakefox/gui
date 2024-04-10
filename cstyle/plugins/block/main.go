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
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]
			parent := s[n.Parent.Properties.Id]

			// If the element is display block and the width is unset then make it 100%

			if n.Style["width"] == "" {
				self.Width, _ = utils.ConvertToPixels("100%", self.EM, parent.Width)
			}
			m := utils.GetMP(*n, "margin")
			self.Width -= (m.Right + m.Left)
			self.Height -= (m.Top + m.Bottom)

			p := utils.GetMP(*n, "padding")
			self.Width += (p.Right + p.Left)
			self.Height += (p.Top + p.Bottom)

			(*state)[n.Properties.Id] = self
		},
	}
}
