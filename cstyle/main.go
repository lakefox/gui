package cstyle

import (
	"fmt"
	"gui/parser"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type CSS struct {
	Width       float32
	Height      float32
	StyleSheets []map[string]map[string]string
}

type Mapped struct {
	Document *html.Node
	StyleMap map[string]map[string]string
}

func (c *CSS) StyleSheet(path string) {
	// Parse the CSS file
	dat, err := os.ReadFile(path)
	check(err)
	styles := parser.ParseCSS(string(dat))

	c.StyleSheets = append(c.StyleSheets, styles)
}

func (c *CSS) StyleTag(css string) {
	styles := parser.ParseCSS(css)
	c.StyleSheets = append(c.StyleSheets, styles)
}

func (c *CSS) Map(doc *html.Node) Mapped {
	styleMap := make(map[string]map[string]string)
	for a := 0; a < len(c.StyleSheets); a++ {
		for key, styles := range c.StyleSheets[a] {
			matching := dom.QuerySelectorAll(doc, key)
			fmt.Printf("%s %s\n", key, matching)
			for _, v := range matching {
				if v.Type == html.ElementNode {
					id := dom.GetAttribute(v, "DOMNODEID")
					if len(id) == 0 {
						id = dom.TagName(v) + fmt.Sprint(rand.Int63())
						dom.SetAttribute(v, "DOMNODEID", id)
					}

					if styleMap[id] == nil {
						styleMap[id] = styles
					} else {
						styleMap[id] = merge(styleMap[id], styles)
					}
				}
			}
		}
	}

	inherit(doc, styleMap)
	size(doc, styleMap, c)

	d := Mapped{
		Document: doc,
		StyleMap: styleMap,
	}
	return d
}

func size(n *html.Node, styleMap map[string]map[string]string, c *CSS) (float32, float32) {
	println("NAME: ", dom.TagName(n))
	id := dom.GetAttribute(n, "DOMNODEID")
	if len(id) == 0 {
		id = dom.TagName(n) + fmt.Sprint(rand.Int63())
		dom.SetAttribute(n, "DOMNODEID", id)
	}
	var width, height float32

	var fixedWidth, fixedHeight bool

	if styleMap[id]["width"] != "" {
		fixedWidth = false
		width, _ = ConvertToPixels(styleMap[id]["width"], c)
	}

	if styleMap[id]["height"] != "" {
		fixedHeight = false
		height, _ = ConvertToPixels(styleMap[id]["height"], c)
	}
	children := dom.Children(n)
	if len(children) > 0 {
		for _, ch := range children {
			w, h := size(ch, styleMap, c)
			if !fixedWidth {
				width = max(w, width)
			}
			if !fixedHeight {
				height = max(h, height)
			}
		}
	} else if styleMap[id]["display"] != "none" {
		text := dom.InnerText(n)
		if len(text) > 0 {
			if styleMap[id]["font-size"] == "" {
				styleMap[id]["font-size"] = "1em"
			}
			fs, _ := ConvertToPixels(styleMap[id]["font-size"], c)
			w, h := getTextBounds(text, fs, width, height)
			fmt.Printf("%f, %f\n", w, h)
			if !fixedWidth {
				width = max(w, width)
			}
			if !fixedHeight {
				height = max(h, height)
			}
		}

	}

	fmt.Printf("RET: %s %f %f\n", dom.TagName(n), width, height)
	return width, height
}

func inherit(n *html.Node, styleMap map[string]map[string]string) {
	if n.Type == html.ElementNode {
		fmt.Println("Element:", n.Data)
		id := dom.GetAttribute(n, "DOMNODEID")
		if len(id) == 0 {
			id = dom.TagName(n) + fmt.Sprint(rand.Int63())
			dom.SetAttribute(n, "DOMNODEID", id)
		}
		pId := dom.GetAttribute(n.Parent, "DOMNODEID")
		if len(pId) > 0 {
			if styleMap[id] == nil {
				styleMap[id] = make(map[string]string)
			}
			if styleMap[pId] == nil {
				styleMap[pId] = make(map[string]string)
			}
			styleMap[id] = xMerge(styleMap[id], styleMap[pId])
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		inherit(c, styleMap)
	}
}

func merge(m1, m2 map[string]string) map[string]string {
	// Create a new map and copy m1 into it
	result := make(map[string]string)
	for k, v := range m1 {
		result[k] = v
	}

	// Merge m2 into the new map
	for k, v := range m2 {
		result[k] = v
	}

	return result
}

func xMerge(m1, m2 map[string]string) map[string]string {
	// Create a new map and copy m1 into it
	result := make(map[string]string)
	for k, v := range m1 {
		result[k] = v
	}

	// Merge m2 into the new map only if the key is not already present
	for k, v := range m2 {
		if result[k] == "" {
			result[k] = v
		}
	}

	return result
}

// ConvertToPixels converts a CSS measurement to pixels.
func ConvertToPixels(value string, c *CSS) (float32, error) {
	// Define conversion factors for different units
	unitFactors := map[string]float32{
		"px": 1,
		"em": 16,    // Assuming 1em = 16px (typical default font size in browsers)
		"pt": 1.33,  // Assuming 1pt = 1.33px (typical conversion)
		"pc": 16.89, // Assuming 1pc = 16.89px (typical conversion)
		"vw": c.Width / 100,
		"vh": c.Height / 100,
	}

	// Extract numeric value and unit using regular expression
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([a-zA-Z]+)$`)
	match := re.FindStringSubmatch(value)

	if len(match) != 3 {
		return 0, fmt.Errorf("invalid input format")
	}

	numericValue, err := (strconv.ParseFloat(match[1], 64))
	numericValue32 := float32(numericValue)
	check(err)

	unit, ok := unitFactors[match[2]]
	if !ok {
		return 0, fmt.Errorf("unsupported unit: %s", match[2])
	}

	return numericValue32 * unit, nil
}

func max(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func getTextBounds(text string, fontSize, width, height float32) (float32, float32) {
	w := float32(len(text) * int(fontSize))
	h := fontSize
	if width > 0 && height > 0 {
		if w > width {
			height = max(height, float32(math.Ceil(float64(w/width)))*h)
		}
		return width, height
	} else {
		return w, h
	}

}
