package main

import (
	"fmt"
	"gui"
	"gui/adapters/raylib"
	"gui/element"
	// "github.com/pkg/profile"
)

// go tool pprof --pdf ./main.go /var/folders/7b/c07zbwkj03nf7cs4vm_0yw1w0000gn/T/profile1893611654/cpu.pprof > file.pdf

func main() {
	// defer profile.Start().Stop() // CPU
	// defer profile.Start(profile.MemProfile).Stop() // Memory
	// defaults read ~/Library/Preferences/.GlobalPreferences.plist

	window := gui.Open("./src/index.html")
	window.Adapter = raylib.Init()
	document := window.Document

	tgt(document.QuerySelector("body"))

	// document.QuerySelector("#editor").AddEventListener("scroll", func(e element.Event) {
	// 	fmt.Println(e.Target.ScrollY, e.Target.TagName)
	// })

	gui.View(&window, 850, 400)
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
