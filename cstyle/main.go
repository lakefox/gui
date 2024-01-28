package cstyle

// package aui/goldie
// https://pkg.go.dev/automated.sh/goldie
// https://pkg.go.dev/automated.sh/aui
// https://pkg.go.dev/automated.sh/oat

import (
	"fmt"
	"gui/color"
	"gui/element"
	"gui/font"
	"gui/parser"
	"gui/utils"
	"hash/fnv"
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

func (c *CSS) CreateDocument(doc *html.Node) element.Node {
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
				Node:   &html.Node{},
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
	return initNodes(&node, *c)
}

func CreateNode(parent element.Node, n *html.Node, slug string) element.Node {
	id := n.Data + slug
	node := element.Node{
		Parent:    &parent,
		TagName:   n.Data,
		InnerText: utils.GetInnerText(n),
		Properties: element.Properties{
			Id:   id,
			Type: n.Type,
			Node: n,
		},
	}
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Split(attr.Val, " ")
			for _, class := range classes {
				node.ClassList.Add(class)
			}
		} else if attr.Key == "id" {
			node.Id = attr.Val
		} else if attr.Key == "contenteditable" && (attr.Val == "" || attr.Val == "true") {
			node.Properties.Editable = true
		} else if attr.Key == "href" {
			node.Href = attr.Val
		} else if attr.Key == "src" {
			node.Src = attr.Val
		} else if attr.Key == "title" {
			node.Title = attr.Val
		}
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

// need to get rid of the .props for the most part all styles should be computed dynamically
// can keep like focusable and stuff that describes the element

// currently the append child does not work due to the props and other stuff not existing so it fails
// moving to a real time style compute would fix that

// :hover is parsed correctly but because the hash func doesn't invalidate it becuase the val
// is updated in the props. change to append :hover to style to create the effect
//							or merge the class with the styles? idk have to think more

func (c *CSS) GetStyles(n element.Node) map[string]string {
	styles := map[string]string{}
	for k, v := range n.Style {
		styles[k] = v
	}
	if n.Parent != nil {
		ps := c.GetStyles(*n.Parent)
		for _, v := range inheritedProps {
			if ps[v] != "" {
				styles[v] = ps[v]
			}
		}

	}

	for _, styleSheet := range c.StyleSheets {
		for selector := range styleSheet {
			// fmt.Println(selector, n.Properties.Id)
			if strings.Contains(selector, "hover") && element.TestSelector(selector, &n) {
				fmt.Println(selector, n.Properties.Id, styleSheet[selector])

			}
			if element.TestSelector(selector, &n) {
				for k, v := range styleSheet[selector] {
					styles[k] = v
				}
			}

		}
	}
	inline := parser.ParseStyleAttribute(n.GetAttribute("style") + ";")
	styles = utils.Merge(styles, inline)
	// add hover and focus css events

	return styles
}

func (c *CSS) Render(doc element.Node) []element.Node {
	return flatten(doc)
}

func (c *CSS) AddPlugin(plugin Plugin) {
	c.Plugins = append(c.Plugins, plugin)
}

func hash(n *element.Node) string {
	// Create a new FNV-1a hash
	hasher := fnv.New32a()

	// Concatenate all values into a single string
	var concatenatedValues string
	for key, val := range n.Style {
		concatenatedValues += key + val
	}
	hasher.Write([]byte(concatenatedValues))

	return string(rune(hasher.Sum32()))
}

func (c *CSS) ComputeNodeStyle(n *element.Node) *element.Node {
	plugins := c.Plugins
	hv := hash(n)
	if n.Properties.Hash != hv {
		// this is kinda a sloppy way to do this but it works ig
		n.Style = c.GetStyles(*n)
		n.Properties.Hash = hv
	}
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
		base := GetPositionOffsetNode(n)
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
		if len(n.InnerText) > 0 {
			innerWidth := width
			innerHeight := height
			genTextNode(n, &innerWidth, &innerHeight)
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
		v.Parent = n
		n.Children[i] = *c.ComputeNodeStyle(&v)
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
			v.Handler(n)
		}
	}

	return n
}

func initNodes(n *element.Node, c CSS) element.Node {
	n = InitNode(n, c)
	for i, ch := range n.Children {
		if ch.Properties.Type == html.ElementNode {
			ch.Parent = n
			cn := initNodes(&ch, c)

			n.Children[i] = cn

		}
	}

	return *n
}

func InitNode(n *element.Node, c CSS) *element.Node {
	n.Style = c.GetStyles(*n)
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
	n.Properties.Text.Font = f
	n.Properties.Text.WordSpacing = int(wordSpacing)
	n.Properties.Text.LetterSpacing = int(letterSpacing)

	n.Properties.Colors = color.Parse(n.Style)
	return n
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

func flatten(n element.Node) []element.Node {
	var nodes []element.Node
	nodes = append(nodes, n)

	children := n.Children
	if len(children) > 0 {
		for _, ch := range children {
			chNodes := flatten(ch)
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

	col, _ := color.Font(n.Style)

	n.Properties.Text.Color = col
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
		*width = utils.Max(*width, float32(font.MeasureLongest(n)))
	} else if n.Style["width"] != "" {
		*width, _ = utils.ConvertToPixels(n.Style["width"], n.Properties.EM, n.Parent.Properties.Width)
	}

	n.Properties.Text.Width = int(*width)
	h := font.Render(n)
	if n.Style["height"] == "" {
		*height = h
	}

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
