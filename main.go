package main

import (
	"fmt"
	rio "gui/document"
	"gui/element"
	// _ "net/http/pprof"
)

func main() {
	// Server for pprof
	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	document := rio.Document{}

	document.Open("./src/app.html", func(doc *element.Node) {
		row := doc.QuerySelector(".row")
		buttons := row.QuerySelectorAll(".button")
		b := *buttons
		for i := range b {
			b[i].AddEventListener("click", func(e element.Event) {
				fmt.Println("Click: ", e.Target.InnerText)
				e.Target.InnerText = "mason"
			})
		}

		editor := doc.QuerySelector("#editor")
		editor.AddEventListener("keypress", func(e element.Event) {
			fmt.Println("key", editor.Value, editor.Properties.Id)
		})
		div := document.CreateElement("div")
		div.InnerText = "test"
		div.ClassList.Add("button")
		row.Children = []element.Node{}
		row.AppendChild(div)
	})

}
