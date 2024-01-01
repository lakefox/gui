package cstyle

// package aui/goldie
// https://pkg.go.dev/automated.sh/goldie
// https://pkg.go.dev/automated.sh/aui

// The font loading needs to be opomised, rn it loads new
// stuff for each one even if they use the same font
// Everything should be one file or at least the rendering pipeline
// Dom needs to be a custom impleamentation for speed and size

import (
	"fmt"
	"gui/color"
	"gui/element"
	"gui/font"
	"gui/parser"
	"gui/utils"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Plugin struct {
	Styles  map[string]string
	Level   int
	Handler func(*element.Node)
}

type CSS struct {
	Width       float32
	Height      float32
	StyleSheets []map[string]map[string]string
	Plugins     []Plugin
	Document    *element.Node
}

type Mapped struct {
	Document *element.Node
	StyleMap map[string]map[string]string
	Render   []element.Node
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

func (c *CSS) CreateDocument(doc *html.Node) {
	id := doc.FirstChild.Data + "0"
	n := doc.FirstChild
	node := element.Node{
		Node: n,
		Parent: &element.Node{
			Id:     "ROOT",
			X:      0,
			Y:      0,
			Width:  c.Width,
			Height: c.Height,
			EM:     16,
			Type:   3,
			Styles: map[string]string{
				"width":  strconv.FormatFloat(float64(c.Width), 'f', -1, 32) + "px",
				"height": strconv.FormatFloat(float64(c.Height), 'f', -1, 32) + "px",
			},
		},
		Id:     id,
		X:      0,
		Y:      0,
		Type:   3,
		Width:  c.Width,
		Height: c.Height,
	}
	i := 0
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			node.Children = append(node.Children, CreateNode(node, child, fmt.Sprint(i)))
			i++
		}
	}
	c.Document = &node
}

func CreateNode(parent element.Node, n *html.Node, slug string) element.Node {
	id := n.Data + slug
	node := element.Node{
		Node:    n,
		Parent:  &parent,
		Type:    n.Type,
		TagName: n.Data,
		Id:      id,
	}
	i := 0
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			node.Children = append(node.Children, CreateNode(node, child, slug+fmt.Sprint(i)))
			i++
		}
	}
	return node
}

// gen id's via a tree so they stay the same
func (c *CSS) Map() Mapped {
	doc := c.Document
	styleMap := make(map[string]map[string]string)
	for a := 0; a < len(c.StyleSheets); a++ {
		for key, styles := range c.StyleSheets[a] {
			matching := doc.QuerySelectorAll(key)
			for _, v := range matching {
				if v.Type == html.ElementNode {
					if styleMap[v.Id] == nil {
						styleMap[v.Id] = styles
					} else {
						styleMap[v.Id] = utils.Merge(styleMap[v.Id], styles)
					}
				}
			}
		}
	}

	// Inherit CSS styles from parent
	inherit(doc, styleMap)
	nodes := initNodes(doc, styleMap)
	node := ComputeNodeStyle(nodes, c.Plugins)
	Print(&node, 0)

	renderLine := flatten(&node)

	d := Mapped{
		Document: &node,
		StyleMap: styleMap,
		Render:   renderLine,
	}
	return d
}

func (c *CSS) AddPlugin(plugin Plugin) {
	c.Plugins = append(c.Plugins, plugin)
}

// make a way of breaking each section out into it's own module so people can add their own.
// this should cover the main parts of html but if some one wants for example drop shadows they
// can make a plug in for it

func ComputeNodeStyle(n element.Node, plugins []Plugin) element.Node {

	styleMap := n.Styles

	if styleMap["display"] == "none" {
		n.X = 0
		n.Y = 0
		n.Width = 0
		n.Height = 0
		return n
	}

	width, height := n.Width, n.Height
	x, y := n.Parent.X, n.Parent.Y

	var top, left, right, bottom bool = false, false, false, false

	if styleMap["position"] == "absolute" {
		base := GetPositionOffsetNode(&n)
		if styleMap["top"] != "" {
			v, _ := utils.ConvertToPixels(styleMap["top"], float32(n.EM), n.Parent.Width)
			y = v + base.Y
			top = true
		}
		if styleMap["left"] != "" {
			v, _ := utils.ConvertToPixels(styleMap["left"], float32(n.EM), n.Parent.Width)
			x = v + base.X
			left = true
		}
		if styleMap["right"] != "" {
			v, _ := utils.ConvertToPixels(styleMap["right"], float32(n.EM), n.Parent.Width)
			x = (base.Width - width) - v
			right = true
		}
		if styleMap["bottom"] != "" {
			v, _ := utils.ConvertToPixels(styleMap["bottom"], float32(n.EM), n.Parent.Width)
			y = (base.Height - height) - v
			bottom = true
		}
	} else {
		for i, v := range n.Parent.Children {
			if v.Id == n.Id {
				if i-1 > 0 {
					sibling := n.Parent.Children[i-1]
					if styleMap["display"] == "inline" {
						if sibling.Styles["display"] == "inline" {
							y = sibling.Y
						} else {
							y = sibling.Y + sibling.Height
						}
					} else {
						y = sibling.Y + sibling.Height
					}
				}
				break
			} else if styleMap["display"] != "inline" {
				y += v.Margin.Top + v.Margin.Bottom + v.Padding.Top + v.Padding.Bottom + v.Height
			}
		}
	}

	// Display modes need to be calculated here

	relPos := !top && !left && !right && !bottom

	if left || relPos {
		x += n.Margin.Left
	}
	if top || relPos {
		y += n.Margin.Top
	}
	if right {
		x -= n.Margin.Right
	}
	if bottom {
		y -= n.Margin.Bottom
	}

	if len(n.Children) == 0 {
		// Confirm text exists
		if len(n.Text.Text) > 0 {
			innerWidth := width
			innerHeight := height
			genTextNode(&n, &innerWidth, &innerHeight)
			width = innerWidth + n.Padding.Left + n.Padding.Right
			height = innerHeight
		}
	}

	n.X = x
	n.Y = y
	n.Width = width
	n.Height = height

	// Sorting the array by the Level field
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Level < plugins[j].Level
	})

	for _, v := range plugins {
		matches := true
		for name, value := range v.Styles {
			if styleMap[name] != value {
				matches = false
			}
		}
		if matches {
			v.Handler(&n)
		}
	}

	// Call children here

	var childYOffset float32
	for i, v := range n.Children {
		v.Parent = &n
		n.Children[i] = ComputeNodeStyle(v, plugins)
		if styleMap["height"] == "" {
			if n.Children[i].Styles["position"] != "absolute" && n.Children[i].Y > childYOffset {
				childYOffset = n.Children[i].Y
				n.Height += n.Children[i].Height
				n.Height += n.Children[i].Margin.Top
				n.Height += n.Children[i].Margin.Bottom
				n.Height += n.Children[i].Padding.Top
				n.Height += n.Children[i].Padding.Bottom
			}

		}
	}

	return n
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
	"display",
}

func inherit(n *element.Node, styleMap map[string]map[string]string) {
	if n.Type == html.ElementNode {
		pId := n.Parent.Id
		if len(pId) > 0 {
			if styleMap[n.Id] == nil {
				styleMap[n.Id] = make(map[string]string)
			}
			if styleMap[pId] == nil {
				styleMap[pId] = make(map[string]string)
			}

			inline := parser.ParseStyleAttribute(n.GetAttribute("style") + ";")
			styleMap[n.Id] = utils.Merge(styleMap[n.Id], inline)
			for _, v := range inheritedProps {
				if styleMap[n.Id][v] == "" && styleMap[pId][v] != "" {
					styleMap[n.Id][v] = styleMap[pId][v]
				}
			}
		}
		utils.SetMP(n.Id, styleMap)
	}

	for _, v := range n.Children {
		inherit(&v, styleMap)
	}
}

func initNodes(n *element.Node, styleMap map[string]map[string]string) element.Node {
	n.Styles = styleMap[n.Id]
	border, err := CompleteBorder(n.Styles)
	if err == nil {
		n.Border = border
	}

	fs, _ := utils.ConvertToPixels(n.Styles["font-size"], n.Parent.EM, n.Parent.Width)
	n.EM = fs

	mt, _ := utils.ConvertToPixels(n.Styles["margin-top"], n.EM, n.Parent.Width)
	mr, _ := utils.ConvertToPixels(n.Styles["margin-right"], n.EM, n.Parent.Width)
	mb, _ := utils.ConvertToPixels(n.Styles["margin-bottom"], n.EM, n.Parent.Width)
	ml, _ := utils.ConvertToPixels(n.Styles["margin-left"], n.EM, n.Parent.Width)
	n.Margin = element.Margin{
		Top:    mt,
		Right:  mr,
		Bottom: mb,
		Left:   ml,
	}

	pt, _ := utils.ConvertToPixels(n.Styles["padding-top"], n.EM, n.Parent.Width)
	pr, _ := utils.ConvertToPixels(n.Styles["padding-right"], n.EM, n.Parent.Width)
	pb, _ := utils.ConvertToPixels(n.Styles["padding-bottom"], n.EM, n.Parent.Width)
	pl, _ := utils.ConvertToPixels(n.Styles["padding-left"], n.EM, n.Parent.Width)
	n.Padding = element.Padding{
		Top:    pt,
		Right:  pr,
		Bottom: pb,
		Left:   pl,
	}

	width, _ := utils.ConvertToPixels(n.Styles["width"], n.EM, n.Parent.Width)
	if n.Styles["min-width"] != "" {
		minWidth, _ := utils.ConvertToPixels(n.Styles["min-width"], n.EM, n.Parent.Width)
		width = utils.Max(width, minWidth)
	}

	if n.Styles["max-width"] != "" {
		maxWidth, _ := utils.ConvertToPixels(n.Styles["max-width"], n.EM, n.Parent.Width)
		width = utils.Min(width, maxWidth)
	}

	height, _ := utils.ConvertToPixels(n.Styles["height"], n.EM, n.Parent.Height)
	if n.Styles["min-height"] != "" {
		minHeight, _ := utils.ConvertToPixels(n.Styles["min-height"], n.EM, n.Parent.Height)
		height = utils.Max(height, minHeight)
	}

	if n.Styles["max-height"] != "" {
		maxHeight, _ := utils.ConvertToPixels(n.Styles["max-height"], n.EM, n.Parent.Height)
		height = utils.Min(height, maxHeight)
	}

	n.Width = width
	n.Height = height

	bold, italic := false, false

	if n.Styles["font-weight"] == "bold" {
		bold = true
	}

	if n.Styles["font-style"] == "italic" {
		italic = true
	}

	f, _ := font.LoadFont(n.Styles["font-family"], int(n.EM), bold, italic)
	letterSpacing, _ := utils.ConvertToPixels(n.Styles["letter-spacing"], n.EM, width)
	wordSpacing, _ := utils.ConvertToPixels(n.Styles["word-spacing"], n.EM, width)
	lineHeight, _ := utils.ConvertToPixels(n.Styles["line-height"], n.EM, width)
	if lineHeight == 0 {
		lineHeight = n.EM + 3
	}

	n.Text.LineHeight = int(lineHeight)
	n.Text.Text = n.InnerText()
	n.Text.Font = f
	n.Text.WordSpacing = int(wordSpacing)
	n.Text.LetterSpacing = int(letterSpacing)

	n.Colors = color.Parse(n.Styles)
	for i, c := range n.Children {
		if c.Type == html.ElementNode {
			c.Parent = n
			cn := initNodes(&c, styleMap)

			n.Children[i] = cn

			if len(n.Children) > 1 {
				cn.PrevSibling = &n.Children[i]
				n.Children[i].NextSibling = &cn
			}
		}
	}

	return *n
}

func GetPositionOffsetNode(n *element.Node) *element.Node {
	pos := n.Styles["position"]

	if pos == "relative" {
		return n
	} else {
		if n.Parent.Node != nil {
			return GetPositionOffsetNode(n.Parent)
		} else {
			return nil
		}
	}
}

func parseBorderShorthand(borderShorthand string) (element.Border, error) {
	// Split the shorthand into components
	borderComponents := strings.Fields(borderShorthand)

	// Ensure there are at least 1 component (width or style or color)
	if len(borderComponents) >= 1 {
		width := "0px" // Default width
		style := "solid"
		borderColor := "#000000" // Default color

		// Extract style and color if available
		if len(borderComponents) >= 1 {
			width = borderComponents[0]
		}

		// Extract style and color if available
		if len(borderComponents) >= 2 {
			style = borderComponents[1]
		}
		if len(borderComponents) >= 3 {
			borderColor = borderComponents[2]
		}

		parsedColor, _ := color.Color(borderColor)

		return element.Border{
			Width:  width,
			Style:  style,
			Color:  parsedColor,
			Radius: "", // Default radius
		}, nil
	}

	return element.Border{}, fmt.Errorf("invalid border shorthand format")
}

func CompleteBorder(cssProperties map[string]string) (element.Border, error) {
	border, err := parseBorderShorthand(cssProperties["border"])
	border.Radius = cssProperties["border-radius"]

	return border, err
}

func Print(n *element.Node, indent int) {
	pre := strings.Repeat("\t", indent)
	fmt.Printf(pre+"%s\n", n.Id)
	fmt.Printf(pre+"-- Parent: %d\n", n.Parent.Id)
	fmt.Printf(pre+"\t-- Width: %f\n", n.Parent.Width)
	fmt.Printf(pre+"\t-- Height: %f\n", n.Parent.Height)
	fmt.Printf(pre + "-- Colors:\n")
	fmt.Printf(pre+"\t-- Font: %f\n", n.Colors.Font)
	fmt.Printf(pre+"\t-- Background: %f\n", n.Colors.Background)
	fmt.Printf(pre+"-- Children: %d\n", len(n.Children))
	fmt.Printf(pre+"-- EM: %f\n", n.EM)
	fmt.Printf(pre+"-- X: %f\n", n.X)
	fmt.Printf(pre+"-- Y: %f\n", n.Y)
	fmt.Printf(pre+"-- Width: %f\n", n.Width)
	fmt.Printf(pre+"-- Height: %f\n", n.Height)
	fmt.Printf(pre+"-- Border: %#v\n", n.Border)
	fmt.Printf(pre+"-- Styles: %#v\n", n.Styles)

	for _, v := range n.Children {
		Print(&v, indent+1)
	}
}

func flatten(n *element.Node) []element.Node {
	var nodes []element.Node
	nodes = append(nodes, *n)

	children := n.Children
	if len(children) > 0 {
		for _, ch := range children {
			chNodes := flatten(&ch)
			nodes = append(nodes, chNodes...)
		}
	}
	return nodes
}

func genTextNode(n *element.Node, width, height *float32) {
	wb := " "

	if n.Styles["word-wrap"] == "break-word" {
		wb = ""
	}

	if n.Styles["text-wrap"] == "wrap" || n.Styles["text-wrap"] == "balance" {
		wb = ""
	}

	letterSpacing, _ := utils.ConvertToPixels(n.Styles["letter-spacing"], n.EM, *width)
	wordSpacing, _ := utils.ConvertToPixels(n.Styles["word-spacing"], n.EM, *width)

	var dt float32

	if n.Styles["text-decoration-thickness"] == "auto" || n.Styles["text-decoration-thickness"] == "" {
		dt = 2
	} else {
		dt, _ = utils.ConvertToPixels(n.Styles["text-decoration-thickness"], n.EM, *width)
	}

	c, _ := color.Font(n.Styles)

	n.Text.Color = c
	n.Text.Align = n.Styles["text-align"]
	n.Text.WordBreak = wb
	n.Text.WordSpacing = int(wordSpacing)
	n.Text.LetterSpacing = int(letterSpacing)
	n.Text.WhiteSpace = n.Styles["white-space"]
	n.Text.DecorationColor = n.Colors.TextDecoration
	n.Text.DecorationThickness = int(dt)
	n.Text.Overlined = n.Styles["text-decoration"] == "overline"
	n.Text.Underlined = n.Styles["text-decoration"] == "underline"
	n.Text.LineThrough = n.Styles["text-decoration"] == "linethrough"
	n.Text.EM = int(n.EM)
	n.Text.Width = int(n.Parent.Width)

	if n.Styles["word-spacing"] == "" {
		n.Text.WordSpacing = font.MeasureSpace(&n.Text)
	}
	if n.Parent.Width != 0 && n.Styles["display"] != "inline" && n.Styles["width"] == "" {
		*width = (n.Parent.Width - n.Padding.Right) - n.Padding.Left
	} else if n.Styles["width"] == "" {
		*width = utils.Max(*width, float32(font.MeasureLongest(&n.Text)))
	} else if n.Styles["width"] != "" {
		*width, _ = utils.ConvertToPixels(n.Styles["width"], n.EM, n.Parent.Width)
	}

	n.Text.Width = int(*width)
	h := font.Render(n)
	if n.Styles["height"] == "" {
		*height = h
	}

}
