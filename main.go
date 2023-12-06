package main

import (
	"fmt"
	"gui/cstyle"
	"gui/document"
)

func main() {
	css := cstyle.CSS{
		Width:  1920,
		Height: 1080,
	}
	css.StyleSheet("./master.css")

	d := document.Write("./src/index.html")

	for _, v := range d.StyleSheets {
		css.StyleSheet(v)
	}

	for _, v := range d.StyleTags {
		css.StyleTag(v)
	}

	// fmt.Printf("%s\n", css.StyleSheets)

	p := css.Map(d.DOM)

	for k, v := range p.StyleMap {
		fmt.Printf("%s\n", k)
		fmt.Printf("%s\n", v)
	}
}
