package block

import (
	"gui/cstyle"
	"gui/element"
)

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"display": "block",
		},
		Level: 0,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]
			// parent := s[n.Parent.Properties.Id]

			// If the element is display block and the width is unset then make it 100%

			// if self.Style["width"] == "" {
			// 	self.Width, _ = utils.ConvertToPixels("100%", self.EM, parent.Width)
			// 	fmt.Println(self.Margin)

			// 	// self.Width -= (self.Padding.Right + self.Padding.Left)
			// 	// self.Height -= (self.Padding.Top + self.Padding.Bottom)
			// }

			// if self.X+self.Width+(self.Border.Width*2) > parent.Width {
			// 	self.Width = parent.Width
			// 	self.Width -= (self.Margin.Right + self.Margin.Left)
			// 	self.Width -= (self.Border.Width * 2)
			// 	self.Height -= (self.Margin.Top + self.Margin.Bottom)
			// }

			(*state)[n.Properties.Id] = self
		},
	}
}
