package ul

import (
	"gui/cstyle"
	"gui/element"
	"strconv"
)

func Init() cstyle.Transformer {
	return cstyle.Transformer{
		Selector: func(n *element.Node) bool {
			return n.TagName == "li"
		},
		Handler: func(n element.Node, c *cstyle.CSS) element.Node {
			if n.Parent.TagName == "ul" {
				dot := n.CreateElement("div")
				dot.Style["background"] = "#000"
				dot.Style["border-radius"] = "100px"
				dot.Style["width"] = "5px"
				dot.Style["height"] = "5px"
				dot.Style["margin-right"] = "10px"

				content := n.CreateElement("div")
				content.InnerText = n.InnerText
				content.Style = n.Style
				content.Style = c.GetStyles(&content)

				n.AppendChild(dot)
				n.AppendChild(content)

			} else if n.Parent.TagName == "ol" {
				dot := n.CreateElement("div")
				// dot.Style["background"] = "#000"
				// dot.Style["border-radius"] = "100px"
				// dot.Style["width"] = "7px"
				// dot.Style["height"] = "7px"
				dot.Style = n.Style
				dot.Style = c.GetStyles(&dot)
				// dot.Style["color"] = "#000"

				for i, v := range n.Parent.Children {
					if v.Properties.Id == n.Properties.Id {
						// var w int
						// for _, v := range c.Fonts {
						// 	w = font.MeasureText(&element.Text{Font: &v}, strconv.Itoa(i+1)+".")
						// 	break
						// }
						dot.InnerText = strconv.Itoa(i+1) + "."
						// fmt.Println(w)
						// dot.Style["margin-right"] = strconv.Itoa(30-w) + "px"
						break
					}
				}

				// !ISSUE: needs to run on parent element so it can know the entire widths of everything, also prob could run styles once?

				content := n.CreateElement("div")
				content.InnerText = n.InnerText
				content.Style = n.Style
				content.Style = c.GetStyles(&content)

				n.AppendChild(dot)
				n.AppendChild(content)

			}
			return n
		},
	}
}
