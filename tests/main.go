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
		fmt.Println(document.QuerySelector("body").InnerHTML)
		// fmt.Println(document.QuerySelector("body").OuterHTML)
		c := document.QuerySelector("h1").Children
		for _, v := range c {
			fmt.Println(v.TagName, v.InnerText)
		}
		// fmt.Println(document.QuerySelector("body").Style)
	})

	// btns := document.QuerySelectorAll(".button")

	// for i := range *btns {
	// 	v := *btns
	// 	v[i].AddEventListener("click", func(e element.Event) {
	// 		fmt.Println(e.Target.InnerText)
	// 		if e.Target.InnerText == "Start" {
	// 			e.Target.InnerText = "Mason"
	// 		}
	// 	})
	// }

	gui.View(&window, 850, 400)

	// input, output := gui.Render(&window, 850, 400)
	// go adapter.View(input, output)
}
