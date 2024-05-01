package main

import (
	"fmt"
	"gui"
	"gui/element"
)

func main() {
	// defaults read ~/Library/Preferences/.GlobalPreferences.plist

	window := gui.Open("./src/index.html")
	// window.AddAdapter(raylib)
	document := window.Document

	document.QuerySelector("body").AddEventListener("click", func(e element.Event) {
		fmt.Println("click")
		// fmt.Println(document.QuerySelector("body").Style)
		document.QuerySelector("body").Style["background"] = "red"
		// fmt.Println(document.QuerySelector("body").Style)
	})
	// test := document.CreateElement("div")
	// test.InnerText = "hellodkljhsa"
	// document.QuerySelector("body").AppendChild(test)

	gui.View(&window, 850, 400)

	// input, output := gui.Render(&window, 850, 400)
	// go adapter.View(input, output)
}
