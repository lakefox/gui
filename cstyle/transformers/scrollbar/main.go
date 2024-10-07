package scrollbar

import (
	"gui/cstyle"
	"gui/element"
	"strconv"
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

			width := "14px"
			if n.Style["scrollbar-width"] == "thin" {
				width = "10px"
			}
			if n.Style["scrollbar-width"] == "none" {
				return n
			}

			splitStr := strings.Split(n.Style["scrollbar-color"], " ")

			// Initialize the variables
			var backgroundColor, thumbColor string

			// Check the length of the split result and assign the values accordingly

			if len(splitStr) >= 2 {
				backgroundColor = splitStr[1]
				thumbColor = splitStr[0]
			} else {
				// backgroundColor = "rgba(0,255,255,1)"
				backgroundColor = "#fafafa"
				// thumbColor = "orange"
				thumbColor = "#c7c7c7"

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
			thumb.Style["top"] = strconv.Itoa(n.ScrollTop) + "px"
			thumb.Style["left"] = "0"
			thumb.Style["width"] = width
			thumb.Style["height"] = "20px"
			thumb.Style["background"] = thumbColor
			thumb.Style["cursor"] = "pointer"
			thumb.Style["z-index"] = "10"
			scrollbar.AppendChild(&thumb)

			n.Style["width"] = "calc(" + n.Style["width"] + "-" + width + ")"
			pr := n.Style["padding-right"]
			if pr == "" {
				if n.Style["padding"] != "" {
					pr = n.Style["padding"]
				}
			}

			if pr != "" {
				n.Style["padding-right"] = "calc(" + pr + "+" + width + ")"
			} else {
				n.Style["padding-right"] = width
			}

			n.AppendChild(&scrollbar)

			return n
		},
	}
}
