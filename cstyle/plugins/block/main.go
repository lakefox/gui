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
				m := utils.GetMP(*n, "margin")
				self.Width -= m.Right + m.Left
			} else {
				p := utils.GetMP(*n, "padding")
				self.Width += p.Right + p.Left
			}
			(*state)[n.Properties.Id] = self
			(*state)[n.Parent.Properties.Id] = parent
		},
	}
}
