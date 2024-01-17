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
			b[i].AddEventListener("click", func(e element.Event) {
				fmt.Println("Click")
			})
		}
	})
}
