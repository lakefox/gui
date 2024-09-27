package scrollbar

import (
	"gui/cstyle"
	"gui/element"
	"gui/utils"
	"strings"
)

func Init() cstyle.Transformer {
	return cstyle.Transformer{
		Selector: func(n *element.Node) bool {
			if n.Style["overflow"] != "" || n.Style["overflow-x"] != "" || n.Style["overflow-y"] != "" {
				return true
			} else {
				return false
			}
		},
		Handler: func(n *element.Node, c *cstyle.CSS) *element.Node {
			// !TODO: Inject grim-scrollbar and grim-scrollbar-thumb elements with stylings added to it
			// + positioning of it should be done by the crop plugin
			// + this is how all :: props should be handled, i think
			// + also crop should hide this if
			// fmt.Println(n.Style, "here")
			if n.Style["position"] == "" {
				n.Style["position"] = "relative"
			}

			width := "20px"
			if n.Style["scrollbar-width"] == "thin" {
				width = "10px"
			}
			if n.Style["scrollbar-width"] == "none" {
				return n
			}

			splitStr := strings.Split(n.Style["scrollbar-width"], " ")

			// Initialize the variables
			var backgroundColor, thumbColor string

			// Check the length of the split result and assign the values accordingly

			if len(splitStr) >= 2 {
				backgroundColor = splitStr[1]
				thumbColor = splitStr[0]
			} else {
				backgroundColor = "rgba(255,0,0,1)"
				// backgroundColor = "#fafafa"
				thumbColor = "orange"
				// thumbColor = "#c7c7c7"

			}

			scrollbar := n.CreateElement("grim-scrollbar")

			scrollbar.Style["position"] = "absolute"
			scrollbar.Style["top"] = "0"
			scrollbar.Style["right"] = "0"
			scrollbar.Style["width"] = width
			scrollbar.Style["height"] = "100%"
			scrollbar.Style["z-index"] = "9"
			scrollbar.Style["background"] = backgroundColor

			thumb := n.CreateElement("grim-thumb")

			thumb.Style["position"] = "absolute"
			thumb.Style["top"] = "10px"
			thumb.Style["left"] = "0"
			thumb.Style["width"] = "18px"
			thumb.Style["height"] = "20px"
			thumb.Style["background"] = thumbColor
			thumb.Style["cursor"] = "pointer"
			thumb.Style["z-index"] = "10"
			scrollbar.AppendChild(&thumb)

			n.AppendChild(&scrollbar)
			return n
		},
	}
}

func AddHTML(n *element.Node) {
	// Head is not renderable
	n.InnerHTML = utils.InnerHTML(n)
	tag, closing := utils.NodeToHTML(n)
	n.OuterHTML = tag + n.InnerHTML + closing
	for i := range n.Children {
		AddHTML(n.Children[i])
	}
}
