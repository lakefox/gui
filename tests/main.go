package main

import (
	"fmt"
	"gui"
	"gui/element"
)

func main() {
	window := gui.Open("./src/index.html")
	document := window.Document

	document.AddEventListener("click", func(e element.Event) {
		fmt.Println("click", e)
		e.Target.Style["background"] = "red"
	})
	test := document.CreateElement("div")
	test.InnerText = "hellodkljhsa"
	document.QuerySelector("body").AppendChild(test)

	gui.View(&window, 850, 400)

	// input, output := gui.Render(&window, 850, 400)
	// go adapter.View(input, output)
}
