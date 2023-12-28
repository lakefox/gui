package main

import (
	"gui/document"
	// _ "net/http/pprof"
)

func main() {
	// Server for pprof
	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	document.Open("./src/app.html")
}
