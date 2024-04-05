package main

import (
	"fmt"
	"gui"
	"gui/element"
)

func main() {
	window := gui.Open("./src/app.html")
	document := window.Document

	document.AddEventListener("click", func(e element.Event) {
		fmt.Println("click", e)
	})

	gui.View(&window, 850, 400)
}
