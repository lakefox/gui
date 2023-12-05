package main

import (
	"fmt"
	"gui/cstyle"
	"gui/document"

	"github.com/go-shiori/dom"
)

func main() {
	css := cstyle.CSS{}
	// css.StyleSheet("./master.css")

	d := document.Parse("./src/index.html")

	for _, v := range d.StyleSheets {
		css.StyleSheet(v)
	}

	for _, v := range d.StyleTags {
		css.StyleTag(v)
	}

	fmt.Printf("%s\n", css.StyleSheets)

	// Example selector: div#test > h1
	selector := "div#test > h1"

	// Use querySelectorAll to find elements that match the selector
	matchingElements := dom.QuerySelectorAll(d.DOM, selector)

	fmt.Printf("%s\n", matchingElements)

	// Print the matching elements
	for _, elem := range matchingElements {
		fmt.Printf("Match: <%s>\n", elem.Data)
	}

}
