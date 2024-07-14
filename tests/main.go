package main

import (
	"fmt"
	"gui"
	"gui/element"
	// "github.com/pkg/profile"
)

// go tool pprof --pdf ./main.go /var/folders/7b/c07zbwkj03nf7cs4vm_0yw1w0000gn/T/profile1893611654/cpu.pprof > file.pdf

func main() {
	// defer profile.Start().Stop() // CPU
	// defer profile.Start(profile.MemProfile).Stop() // Memory
	// defaults read ~/Library/Preferences/.GlobalPreferences.plist

	window := gui.Open("./src/index.html")
	// window.AddAdapter(raylib)
	document := window.Document

	tgt(document.QuerySelector("body"))

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
func tgt(e *element.Node) {
	// events need to be transfered to broke out elements
	e.AddEventListener("click", func(e element.Event) {
		// fmt.Println(document.QuerySelector("body").InnerHTML)
		fmt.Println(e.Target.TagName)
		e.Target.Style["background"] = "red"
	})

	for i := range e.Children {

		tgt(e.Children[i])
	}
}
