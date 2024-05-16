package text

import (
	"gui/cstyle"
	"gui/element"
	"gui/utils"
	"strings"
)

func Init() cstyle.Transformer {
	return cstyle.Transformer{
		Selector: func(n *element.Node) bool {
			if !utils.ChildrenHaveText(n) && len(n.InnerText) > 0 {
				// Confirm text exists
				words := strings.Split(strings.TrimSpace(n.InnerText), " ")
				return len(words) != 1
			} else {
				return false
			}
			// return true
		},
		Handler: func(n element.Node) element.Node {
			if utils.IsParent(n, "head") {
				return n
			}
			words := strings.Split(strings.TrimSpace(n.InnerText), " ")
			n.InnerText = words[0]
			// fmt.Println("##########")
			// fmt.Println(n.TagName)
			// !ISSUE: issue is here don't know why
			for i := 1; i < len(words); i++ {
				el := n
				el.InnerText = words[i]
				n.Parent.InsertAfter(el, n)
				// fmt.Println("inject", el.Properties.Id)
			}
			return n
		},
	}
}
