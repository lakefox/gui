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
}

func (c *CSS) StyleSheet(path string) {
	// Parse the CSS file
	dat, err := os.ReadFile(path)
	utils.Check(err)
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
		Parent: &element.Node{
			Properties: element.Properties{
				Id:     "ROOT",
				X:      0,
				Y:      0,
				Width:  c.Width,
				Height: c.Height,
				EM:     16,
				Type:   3,
			},

			Style: map[string]string{
				"width":  strconv.FormatFloat(float64(c.Width), 'f', -1, 32) + "px",
				"height": strconv.FormatFloat(float64(c.Height), 'f', -1, 32) + "px",
			},
		},
		Properties: element.Properties{
			Node:   n,
			Id:     id,
			X:      0,
			Y:      0,
			Type:   3,
			Width:  c.Width,
			Height: c.Height,
		},
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
		Parent:  &parent,
		TagName: n.Data,
		Properties: element.Properties{
			Id:   id,
			Type: n.Type,
			Node: n,
		},
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
			for _, v := range *matching {
				if v.Properties.Type == html.ElementNode {
					if styleMap[v.Properties.Id] == nil {
						styleMap[v.Properties.Id] = styles
					} else {
						styleMap[v.Properties.Id] = utils.Merge(styleMap[v.Properties.Id], styles)
					}
				}
			}
		}
	}

	// Inherit CSS styles from parent
	inherit(doc, styleMap)
	nodes := initNodes(doc, styleMap)
	node := ComputeNodeStyle(nodes, c.Plugins)
	// Print(&node, 0)

	d := Mapped{
		Document: &node,
		StyleMap: styleMap,
	}
	return d
}

func (m *Mapped) Render() []element.Node {
	return flatten(m.Document)
}

func (c *CSS) AddPlugin(plugin Plugin) {
	c.Plugins = append(c.Plugins, plugin)
}

// make a way of breaking each section out into it's own module so people can add their own.
// this should cover the main parts of html but if some one wants for example drop shadows they
// can make a plug in for it

func ComputeNodeStyle(n element.Node, plugins []Plugin) element.Node {

	styleMap := n.Style

	if styleMap["display"] == "none" {
		n.Properties.X = 0
		n.Properties.Y = 0
		n.Properties.Width = 0
		n.Properties.Height = 0
		return n
	}

	width, height := n.Properties.Width, n.Properties.Height
	x, y := n.Parent.Properties.X, n.Parent.Properties.Y

	var top, left, right, bottom bool = false, false, false, false

	if styleMap["position"] == "absolute" {
		base := GetPositionOffsetNode(&n)
		if styleMap["top"] != "" {
			v, _ := utils.ConvertToPixels(styleMap["top"], float32(n.Properties.EM), n.Parent.Properties.Width)
			y = v + base.Properties.Y
			top = true
		}
		if styleMap["left"] != "" {
			v, _ := utils.ConvertToPixels(styleMap["left"], float32(n.Properties.EM), n.Parent.Properties.Width)
			x = v + base.Properties.X
			left = true
		}
		if styleMap["right"] != "" {
			v, _ := utils.ConvertToPixels(styleMap["right"], float32(n.Properties.EM), n.Parent.Properties.Width)
			x = (base.Properties.Width - width) - v
			right = true
		}
		if styleMap["bottom"] != "" {
			v, _ := utils.ConvertToPixels(styleMap["bottom"], float32(n.Properties.EM), n.Parent.Properties.Width)
			y = (base.Properties.Height - height) - v
			bottom = true
		}
	} else {
		for i, v := range n.Parent.Children {
			if v.Properties.Id == n.Properties.Id {
				if i-1 > 0 {
					sibling := n.Parent.Children[i-1]
					if styleMap["display"] == "inline" {
						if sibling.Style["display"] == "inline" {
							y = sibling.Properties.Y
						} else {
							y = sibling.Properties.Y + sibling.Properties.Height
						}
					} else {
						y = sibling.Properties.Y + sibling.Properties.Height
					}
				}
				break
			} else if styleMap["display"] != "inline" {
				y += v.Properties.Margin.Top + v.Properties.Margin.Bottom + v.Properties.Padding.Top + v.Properties.Padding.Bottom + v.Properties.Height
			}
		}
	}

	// Display modes need to be calculated here

	relPos := !top && !left && !right && !bottom

	if left || relPos {
		x += n.Properties.Margin.Left
	}
	if top || relPos {
		y += n.Properties.Margin.Top
	}
	if right {
		x -= n.Properties.Margin.Right
	}
	if bottom {
		y -= n.Properties.Margin.Bottom
	}

	if len(n.Children) == 0 {
		// Confirm text exists
		if len(n.Properties.Text.Text) > 0 {
			innerWidth := width
			innerHeight := height
			genTextNode(&n, &innerWidth, &innerHeight)
			width = innerWidth + n.Properties.Padding.Left + n.Properties.Padding.Right
			height = innerHeight
		}
	}

	n.Properties.X = x
	n.Properties.Y = y
	n.Properties.Width = width
	n.Properties.Height = height

	// Call children here

	var childYOffset float32
	for i, v := range n.Children {
		v.Parent = &n
		n.Children[i] = ComputeNodeStyle(v, plugins)
		if styleMap["height"] == "" {
			if n.Children[i].Style["position"] != "absolute" && n.Children[i].Properties.Y > childYOffset {
				childYOffset = n.Children[i].Properties.Y
				n.Properties.Height += n.Children[i].Properties.Height
				n.Properties.Height += n.Children[i].Properties.Margin.Top
				n.Properties.Height += n.Children[i].Properties.Margin.Bottom
				n.Properties.Height += n.Children[i].Properties.Padding.Top
				n.Properties.Height += n.Children[i].Properties.Padding.Bottom
			}

		}
	}

	// Sorting the array by the Level field
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Level < plugins[j].Level
	})

	for _, v := range plugins {
		matches := true
		for name, value := range v.Styles {
			if styleMap[name] != value && !(value == "*") {
				matches = false
			}
		}
		if matches {
			v.Handler(&n)
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
	if n.Properties.Type == html.ElementNode {
		pId := n.Parent.Properties.Id
		if len(pId) > 0 {
			if styleMap[n.Properties.Id] == nil {
				styleMap[n.Properties.Id] = make(map[string]string)
			}
			if styleMap[pId] == nil {
				styleMap[pId] = make(map[string]string)
			}

			inline := parser.ParseStyleAttribute(n.GetAttribute("style") + ";")
			styleMap[n.Properties.Id] = utils.Merge(styleMap[n.Properties.Id], inline)
			for _, v := range inheritedProps {
				if styleMap[n.Properties.Id][v] == "" && styleMap[pId][v] != "" {
					styleMap[n.Properties.Id][v] = styleMap[pId][v]
				}
			}
		}
		utils.SetMP(n.Properties.Id, styleMap)
	}

	for _, v := range n.Children {
		inherit(&v, styleMap)
	}
}

func initNodes(n *element.Node, styleMap map[string]map[string]string) element.Node {
	n.Style = styleMap[n.Properties.Id]
	border, err := CompleteBorder(n.Style)
	if err == nil {
		n.Properties.Border = border
	}

	fs, _ := utils.ConvertToPixels(n.Style["font-size"], n.Parent.Properties.EM, n.Parent.Properties.Width)
	n.Properties.EM = fs

	mt, _ := utils.ConvertToPixels(n.Style["margin-top"], n.Properties.EM, n.Parent.Properties.Width)
	mr, _ := utils.ConvertToPixels(n.Style["margin-right"], n.Properties.EM, n.Parent.Properties.Width)
	mb, _ := utils.ConvertToPixels(n.Style["margin-bottom"], n.Properties.EM, n.Parent.Properties.Width)
	ml, _ := utils.ConvertToPixels(n.Style["margin-left"], n.Properties.EM, n.Parent.Properties.Width)
	n.Properties.Margin = element.Margin{
		Top:    mt,
		Right:  mr,
		Bottom: mb,
		Left:   ml,
	}

	pt, _ := utils.ConvertToPixels(n.Style["padding-top"], n.Properties.EM, n.Parent.Properties.Width)
	pr, _ := utils.ConvertToPixels(n.Style["padding-right"], n.Properties.EM, n.Parent.Properties.Width)
	pb, _ := utils.ConvertToPixels(n.Style["padding-bottom"], n.Properties.EM, n.Parent.Properties.Width)
	pl, _ := utils.ConvertToPixels(n.Style["padding-left"], n.Properties.EM, n.Parent.Properties.Width)
	n.Properties.Padding = element.Padding{
		Top:    pt,
		Right:  pr,
		Bottom: pb,
		Left:   pl,
	}

	width, _ := utils.ConvertToPixels(n.Style["width"], n.Properties.EM, n.Parent.Properties.Width)
	if n.Style["min-width"] != "" {
		minWidth, _ := utils.ConvertToPixels(n.Style["min-width"], n.Properties.EM, n.Parent.Properties.Width)
		width = utils.Max(width, minWidth)
	}

	if n.Style["max-width"] != "" {
		maxWidth, _ := utils.ConvertToPixels(n.Style["max-width"], n.Properties.EM, n.Parent.Properties.Width)
		width = utils.Min(width, maxWidth)
	}

	height, _ := utils.ConvertToPixels(n.Style["height"], n.Properties.EM, n.Parent.Properties.Height)
	if n.Style["min-height"] != "" {
		minHeight, _ := utils.ConvertToPixels(n.Style["min-height"], n.Properties.EM, n.Parent.Properties.Height)
		height = utils.Max(height, minHeight)
	}

	if n.Style["max-height"] != "" {
		maxHeight, _ := utils.ConvertToPixels(n.Style["max-height"], n.Properties.EM, n.Parent.Properties.Height)
		height = utils.Min(height, maxHeight)
	}

	if n.Style["margin"] == "auto" && n.Style["margin-left"] == "" && n.Style["margin-right"] == "" {
		n.Properties.Margin.Left = utils.Max((n.Parent.Properties.Width-width)/2, 0)
		n.Properties.Margin.Right = utils.Max((n.Parent.Properties.Width-width)/2, 0)
	}

	n.Properties.Width = width
	n.Properties.Height = height

	bold, italic := false, false

	if n.Style["font-weight"] == "bold" {
		bold = true
	}

	if n.Style["font-style"] == "italic" {
		italic = true
	}

	f, _ := font.LoadFont(n.Style["font-family"], int(n.Properties.EM), bold, italic)
	letterSpacing, _ := utils.ConvertToPixels(n.Style["letter-spacing"], n.Properties.EM, width)
	wordSpacing, _ := utils.ConvertToPixels(n.Style["word-spacing"], n.Properties.EM, width)
	lineHeight, _ := utils.ConvertToPixels(n.Style["line-height"], n.Properties.EM, width)
	if lineHeight == 0 {
		lineHeight = n.Properties.EM + 3
	}

	n.Properties.Text.LineHeight = int(lineHeight)
	n.Properties.Text.Text = n.InnerText()
	n.Properties.Text.Font = f
	n.Properties.Text.WordSpacing = int(wordSpacing)
	n.Properties.Text.LetterSpacing = int(letterSpacing)

	n.Properties.Colors = color.Parse(n.Style)
	for i, c := range n.Children {
		if c.Properties.Type == html.ElementNode {
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
	pos := n.Style["position"]

	if pos == "relative" {
		return n
	} else {
		if n.Parent.Properties.Node != nil {
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
	fmt.Printf(pre+"%s\n", n.Properties.Id)
	fmt.Printf(pre+"-- Parent: %d\n", n.Parent.Properties.Id)
	fmt.Printf(pre+"\t-- Width: %f\n", n.Parent.Properties.Width)
	fmt.Printf(pre+"\t-- Height: %f\n", n.Parent.Properties.Height)
	fmt.Printf(pre + "-- Colors:\n")
	fmt.Printf(pre+"\t-- Font: %f\n", n.Properties.Colors.Font)
	fmt.Printf(pre+"\t-- Background: %f\n", n.Properties.Colors.Background)
	fmt.Printf(pre+"-- Children: %d\n", len(n.Children))
	fmt.Printf(pre+"-- EM: %f\n", n.Properties.EM)
	fmt.Printf(pre+"-- X: %f\n", n.Properties.X)
	fmt.Printf(pre+"-- Y: %f\n", n.Properties.Y)
	fmt.Printf(pre+"-- Width: %f\n", n.Properties.Width)
	fmt.Printf(pre+"-- Height: %f\n", n.Properties.Height)
	fmt.Printf(pre+"-- Border: %#v\n", n.Properties.Border)
	fmt.Printf(pre+"-- Styles: %#v\n", n.Style)

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

	if n.Style["word-wrap"] == "break-word" {
		wb = ""
	}

	if n.Style["text-wrap"] == "wrap" || n.Style["text-wrap"] == "balance" {
		wb = ""
	}

	letterSpacing, _ := utils.ConvertToPixels(n.Style["letter-spacing"], n.Properties.EM, *width)
	wordSpacing, _ := utils.ConvertToPixels(n.Style["word-spacing"], n.Properties.EM, *width)

	var dt float32

	if n.Style["text-decoration-thickness"] == "auto" || n.Style["text-decoration-thickness"] == "" {
		dt = 2
	} else {
		dt, _ = utils.ConvertToPixels(n.Style["text-decoration-thickness"], n.Properties.EM, *width)
	}

	c, _ := color.Font(n.Style)

	n.Properties.Text.Color = c
	n.Properties.Text.Align = n.Style["text-align"]
	n.Properties.Text.WordBreak = wb
	n.Properties.Text.WordSpacing = int(wordSpacing)
	n.Properties.Text.LetterSpacing = int(letterSpacing)
	n.Properties.Text.WhiteSpace = n.Style["white-space"]
	n.Properties.Text.DecorationColor = n.Properties.Colors.TextDecoration
	n.Properties.Text.DecorationThickness = int(dt)
	n.Properties.Text.Overlined = n.Style["text-decoration"] == "overline"
	n.Properties.Text.Underlined = n.Style["text-decoration"] == "underline"
	n.Properties.Text.LineThrough = n.Style["text-decoration"] == "linethrough"
	n.Properties.Text.EM = int(n.Properties.EM)
	n.Properties.Text.Width = int(n.Parent.Properties.Width)

	if n.Style["word-spacing"] == "" {
		n.Properties.Text.WordSpacing = font.MeasureSpace(&n.Properties.Text)
	}
	if n.Parent.Properties.Width != 0 && n.Style["display"] != "inline" && n.Style["width"] == "" {
		*width = (n.Parent.Properties.Width - n.Properties.Padding.Right) - n.Properties.Padding.Left
	} else if n.Style["width"] == "" {
		*width = utils.Max(*width, float32(font.MeasureLongest(&n.Properties.Text)))
	} else if n.Style["width"] != "" {
		*width, _ = utils.ConvertToPixels(n.Style["width"], n.Properties.EM, n.Parent.Properties.Width)
	}

	n.Properties.Text.Width = int(*width)
	h := font.Render(n)
	if n.Style["height"] == "" {
		*height = h
	}

}
