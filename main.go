package main

import (
	"fmt"
	"gui/document"
	"gui/element"
	// _ "net/http/pprof"
)

func main() {
	// Server for pprof
	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	document.Open("./src/app.html", func(doc *element.Node) {
		row := doc.QuerySelector(".row")
		buttons := row.QuerySelectorAll(".button")
		b := *buttons
		for i := range b {
			// b[i].InnerText = "mason"
			b[i].AddEventListener("click", func(e element.Event) {
				fmt.Println("Click: ", e.Target.InnerText)
			})
			b[i].AddEventListener("mouseenter", func(e element.Event) {
				fmt.Println("MOUSE ENTER")
			})
			b[i].AddEventListener("mouseleave", func(e element.Event) {
				fmt.Println("MOUSE LEAVE")
			})
			b[i].AddEventListener("mouseover", func(e element.Event) {
				fmt.Println("MOUSE OVER")
			})

			b[i].AddEventListener("mousemove", func(e element.Event) {
				fmt.Println("MOUSE POSITION: ", e.X, e.Y)
			})
		}

		doc.AddEventListener("mousedown", func(e element.Event) {
			fmt.Println("MOUSE DOWN")
		})

		doc.AddEventListener("mouseup", func(e element.Event) {
			fmt.Println("MOUSE UP")
		})

		doc.AddEventListener("click", func(e element.Event) {
			fmt.Println("CLICK")
		})

		doc.AddEventListener("scroll", func(e element.Event) {
			fmt.Println("Y: ", e.Target.ScrollY)
		})
		editor := doc.QuerySelector("#editor")
		editor.AddEventListener("keypress", func(e element.Event) {
			fmt.Println("key", editor.Value, editor.Properties.Id)
		})
	})
}
