package flex

import (
	"fmt"
	"gui/cstyle"
	"gui/element"
)

func Plugin() cstyle.Plugin {
	return cstyle.Plugin{
		Styles: map[string]string{
			"display": "inline",
		},
		Level: 0,
		Handler: func(n *element.Node) {
			fmt.Println("hi")
		},
	}
}
