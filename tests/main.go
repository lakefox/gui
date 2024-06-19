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

	window := gui.Open("./src/flexcol.html")
	// window.AddAdapter(raylib)
	document := window.Document

	document.QuerySelector("body").AddEventListener("click", func(e element.Event) {
		fmt.Println("click")
		fmt.Println(document.QuerySelector("body").InnerHTML)
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
