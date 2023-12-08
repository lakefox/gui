package cstyle

import (
	"fmt"
	"gui/parser"
	"math/rand"
	"os"

	"gui/utils"

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
	Render   []Node
}

type Node struct {
	Node *html.Node
	Id   string
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
						styleMap[id] = utils.Merge(styleMap[id], styles)
					}
				}
			}
		}
	}

	// Inherit CSS styles from parent
	inherit(doc, styleMap)
	// Calculate the width and height
	fmt.Printf("123 %f %f\n", c.Width, c.Height)
	size(doc, styleMap, c.Width, c.Height)
	// Calculate the X and Y values
	position(doc, styleMap, 0, 0, c.Width, c.Height, c.Width, c.Height)

	renderLine := flatten(doc)

	d := Mapped{
		Document: doc,
		StyleMap: styleMap,
		Render:   renderLine,
	}
	return d
}

func flatten(n *html.Node) []Node {
	var nodes []Node
	id := dom.GetAttribute(n, "DOMNODEID")
	nodes = append(nodes, Node{
		Node: n,
		Id:   id,
	})

	children := dom.Children(n)
	if len(children) > 0 {
		for _, ch := range children {
			chNodes := flatten(ch)
			nodes = append(nodes, chNodes...)
		}
	}
	return nodes
}

func position(n *html.Node, styleMap map[string]map[string]string, x1, y1, x2, y2, windowWidth, windowHeight float32) (float32, float32, float32, float32) {
	id := dom.GetAttribute(n, "DOMNODEID")
	if len(id) == 0 {
		id = dom.TagName(n) + fmt.Sprint(rand.Int63())
		dom.SetAttribute(n, "DOMNODEID", id)
	}

	width, _ := utils.ConvertToPixels(styleMap[id]["width"], windowWidth)
	height, _ := utils.ConvertToPixels(styleMap[id]["height"], windowHeight)

	x2 = width
	y2 = height

	if styleMap[id]["margin-left"] != "" {
		v, _ := utils.ConvertToPixels(styleMap[id]["margin-left"], windowWidth)
		x1 += v
	}
	if styleMap[id]["margin-top"] != "" {
		v, _ := utils.ConvertToPixels(styleMap[id]["margin-top"], windowHeight)
		y1 += v
	}

	if styleMap[id]["margin-right"] != "" {
		v, _ := utils.ConvertToPixels(styleMap[id]["margin-left"], windowWidth)
		x2 += v
	}
	if styleMap[id]["margin-bottom"] != "" {
		v, _ := utils.ConvertToPixels(styleMap[id]["margin-top"], windowHeight)
		y2 += v
	}

	children := dom.Children(n)

	if len(children) > 0 {
		for _, ch := range children {
			_, b, _, d := position(ch, styleMap, x1, y1, x2, y2, width, height)
			y1 += b + d
		}
	}
	if styleMap[id] == nil {
		styleMap[id] = make(map[string]string)
	}
	styleMap[id]["x"] = fmt.Sprintf("%g", x1)
	styleMap[id]["y"] = fmt.Sprintf("%g", y1)
	return x1, y1, x2, y2
}

func size(n *html.Node, styleMap map[string]map[string]string, windowWidth, windowHeight float32) (float32, float32) {
	fmt.Printf("%f %f\n", windowWidth, windowHeight)
	id := dom.GetAttribute(n, "DOMNODEID")
	if len(id) == 0 {
		id = dom.TagName(n) + fmt.Sprint(rand.Int63())
		dom.SetAttribute(n, "DOMNODEID", id)
	}
	var width, height float32

	if styleMap[id]["width"] != "" {
		width, _ = utils.ConvertToPixels(styleMap[id]["width"], windowWidth)
		fmt.Printf("%f %f %s %s\n", width, windowWidth, dom.TagName(n), styleMap[id]["width"])
		fmt.Printf("%s\n", styleMap[id])
		t, _ := utils.ConvertToPixels("50%", 100)
		fmt.Printf("%s\n", t)
	}

	if styleMap[id]["height"] != "" {
		height, _ = utils.ConvertToPixels(styleMap[id]["height"], windowHeight)
	}

	children := dom.Children(n)
	if len(children) > 0 {
		for _, ch := range children {
			if width == 0 {
				width = windowWidth
			}
			if height == 0 {
				height = windowHeight
			}
			w, h := size(ch, styleMap, width, height)

			width = utils.Max(w, width)

			height += h

		}
	} else if styleMap[id]["display"] != "none" {
		text := dom.InnerText(n)
		if len(text) > 0 {
			if styleMap[id]["font-size"] == "" {
				styleMap[id]["font-size"] = "1em"
			}
			fs, _ := utils.ConvertToPixels(styleMap[id]["font-size"], width)
			w, h := utils.GetTextBounds(text, fs, width, height)

			width = w

			height = h

		}

	}

	width, height = utils.AddMarginAndPadding(styleMap, id, width, height)

	if styleMap[id] == nil {
		styleMap[id] = make(map[string]string)
	}
	styleMap[id]["width"] = fmt.Sprintf("%g", width)
	styleMap[id]["height"] = fmt.Sprintf("%g", height)
	return width, height
}

func inherit(n *html.Node, styleMap map[string]map[string]string) {
	if n.Type == html.ElementNode {
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
			styleMap[id] = utils.ExMerge(styleMap[id], styleMap[pId])
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		inherit(c, styleMap)
	}
}
