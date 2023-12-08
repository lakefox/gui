package cstyle

import (
	"fmt"
	"gui/parser"
	"math/rand"
	"os"
	"strconv"

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
	println("inherit")
	inherit(doc, styleMap)
	// Calculate the width and height
	println("size")
	size(doc, styleMap, c.Width, c.Height)
	// Calculate the X and Y values
	println("position")
	position(doc, styleMap, 0, 0, c.Width, c.Height)

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

func position(n *html.Node, styleMap map[string]map[string]string, x1, y1, windowWidth, windowHeight float32) (float32, float32) {
	id := dom.GetAttribute(n, "DOMNODEID")
	println(dom.TagName(n))
	if len(id) == 0 {
		id = dom.TagName(n) + fmt.Sprint(rand.Int63())
		dom.SetAttribute(n, "DOMNODEID", id)
	}

	fs := utils.GetFontSize(styleMap[id])

	rawWidth, _ := strconv.ParseFloat(styleMap[id]["width"], 32)
	rawHeight, _ := strconv.ParseFloat(styleMap[id]["height"], 32)
	width := float32(rawWidth)
	height := float32(rawHeight)

	x2 := x1 + width
	y2 := y1 + height

	var btmOS float32 = 0

	if styleMap[id]["margin-left"] != "" {
		v, _ := utils.ConvertToPixels(styleMap[id]["margin-left"], float32(fs), windowWidth)
		x1 += v
		x2 += v
	}
	if styleMap[id]["margin-top"] != "" {
		v, _ := utils.ConvertToPixels(styleMap[id]["margin-top"], float32(fs), windowHeight)
		y1 += v
		y2 += v
		btmOS += v
	}

	if styleMap[id]["margin-bottom"] != "" {
		v, _ := utils.ConvertToPixels(styleMap[id]["margin-top"], float32(fs), windowHeight)
		btmOS += v
	}

	children := dom.Children(n)
	oY := btmOS
	if len(children) > 0 {
		for _, ch := range children {
			_, h := position(ch, styleMap, x1, y1+oY, width, height)

			oY += h
		}
	}
	if styleMap[id] == nil {
		styleMap[id] = make(map[string]string)
	}
	styleMap[id]["x"] = fmt.Sprintf("%g", x1)
	styleMap[id]["y"] = fmt.Sprintf("%g", y1)
	return x2 - x1, (y2 + btmOS) - y1
}

func size(n *html.Node, styleMap map[string]map[string]string, windowWidth, windowHeight float32) (float32, float32) {
	id := dom.GetAttribute(n, "DOMNODEID")
	println(dom.TagName(n))
	if len(id) == 0 {
		id = dom.TagName(n) + fmt.Sprint(rand.Int63())
		dom.SetAttribute(n, "DOMNODEID", id)
	}

	fs := utils.GetFontSize(styleMap[id])

	var width, height float32

	if styleMap[id]["width"] != "" {
		width, _ = utils.ConvertToPixels(styleMap[id]["width"], float32(fs), windowWidth)
	}

	if styleMap[id]["height"] != "" {
		height, _ = utils.ConvertToPixels(styleMap[id]["height"], float32(fs), windowHeight)
	}

	children := dom.Children(n)
	if len(children) > 0 {
		for _, ch := range children {
			var wW, wH float32 = width, height
			if n.Type != html.ElementNode {
				wW = windowWidth
				wH = windowHeight
			}
			w, h := size(ch, styleMap, wW, wH)

			width = utils.Max(width, w)

			height += h

		}
	} else if styleMap[id]["display"] != "none" {
		text := dom.InnerText(n)
		if len(text) > 0 {
			if styleMap[id]["font-size"] == "" {
				styleMap[id]["font-size"] = "1em"
			}
			fs2, _ := utils.ConvertToPixels(styleMap[id]["font-size"], fs, width)

			_, h := utils.GetTextBounds(text, fs2, width, height)

			height = h

		}

	}
	var (
		wMarginWidth  float32
		wMarginHeight float32
	)

	utils.SetMP(id, styleMap)

	width, height, wMarginWidth, wMarginHeight = utils.AddMarginAndPadding(styleMap, id, width, height)

	if styleMap[id] == nil {
		styleMap[id] = make(map[string]string)
	}

	_, right, _, left := utils.GetMarginOffset(n, styleMap, windowWidth, windowHeight)

	styleMap[id]["width"] = fmt.Sprintf("%g", width-(left+right))
	styleMap[id]["height"] = fmt.Sprintf("%g", height)
	return wMarginWidth, wMarginHeight
}

var inheritedProps = []string{
	"color",
	"cursor",
	"font",
	"font-family",
	"font-size",
	"font-style",
	"font-weight",
	"letter-spacing",
	"line-height",
	"text-align",
	"text-indent",
	"text-justify",
	"text-shadow",
	"text-transform",
	"visibility",
	"word-spacing",
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

			for _, v := range inheritedProps {
				if styleMap[id][v] == "" && styleMap[pId][v] != "" {
					styleMap[id][v] = styleMap[pId][v]
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		inherit(c, styleMap)
	}
}
