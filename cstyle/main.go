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
	"math/rand"
	"os"
	"strconv"
	"strings"

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
	fId := dom.GetAttribute(doc.FirstChild, "DOMNODEID")
	node := element.Node{
		Node: doc.FirstChild,
		Parent: &element.Node{
			Id:     "ROOT",
			X:      0,
			Y:      0,
			Width:  c.Width,
			Height: c.Height,
			EM:     16,
			Styles: map[string]string{
				"width":  strconv.FormatFloat(float64(c.Width), 'f', -1, 32) + "px",
				"height": strconv.FormatFloat(float64(c.Height), 'f', -1, 32) + "px",
			},
		},
		Id:     fId,
		X:      0,
		Y:      0,
		Width:  c.Width,
		Height: c.Height,
		Styles: styleMap[fId],
	}
	initNodes(&node, styleMap)

	node = ComputeNodeStyle(node)

	// Print(&node, 0)

	renderLine := flatten(&node)

	d := Mapped{
		Document: &node,
		StyleMap: styleMap,
		Render:   renderLine,
	}
	return d
}

// make a way of breaking each section out into it's own module so people can add their own.
// this should cover the main parts of html but if some one wants for example drop shadows they
// can make a plug in for it

func ComputeNodeStyle(n element.Node) element.Node {

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

	if styleMap["display"] == "block" {
		// If the element is display block and the width is unset then make it 100%
		if styleMap["width"] == "" {
			width, _ = utils.ConvertToPixels("100%", n.EM, n.Parent.Width)
			width -= n.Margin.Right + n.Margin.Left
		} else {
			width += n.Padding.Right + n.Padding.Left
		}
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

	if styleMap["display"] == "inline" {
		copyOfX := x
		xAcume := float32(0)
		for _, v := range n.Parent.Children {
			if v.Id == n.Id {
				break
			} else if v.Styles["display"] == "inline" {
				fmt.Println(x+xAcume+n.Width, n.Parent.Width, x, y)
				if x+xAcume+n.Width > n.Parent.Width {
					y += float32(n.Text.LineHeight)
					n.Parent.Height += float32(n.Text.LineHeight)
					x = copyOfX
					xAcume = 0
				} else {
					x += v.Width
				}
			} else {
				x = copyOfX
			}
			xAcume += v.X + v.Width
		}

	}
	fmt.Println(x, y)

	n.X = x
	n.Y = y
	n.Width = width
	n.Height = height

	// Call children here

	var childYOffset float32

	for i, v := range n.Children {
		v.Parent = &n
		n.Children[i] = ComputeNodeStyle(v)
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

	if styleMap["text-align"] == "center" || styleMap["text-align"] == "right" {
		pos := float32(1)
		if styleMap["text-align"] == "center" {
			pos = 2
		}
		offset := n.Width / pos
		for i, v := range n.Children {
			if v.Styles["display"] == "inline" {
				offset -= v.Width / pos
			} else {
				for j := i - 1; j >= 0; j-- {
					if n.Children[j].Styles["display"] != "inline" {
						break
					} else {
						n.Children[j].X += offset
					}
				}
				offset = n.Width / pos
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

			inline := parser.ParseStyleAttribute(dom.GetAttribute(n, "style") + ";")
			styleMap[id] = utils.Merge(styleMap[id], inline)

			for _, v := range inheritedProps {
				if styleMap[id][v] == "" && styleMap[pId][v] != "" {
					styleMap[id][v] = styleMap[pId][v]
				}
			}
		}
		utils.SetMP(id, styleMap)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		inherit(c, styleMap)
	}
}

func initNodes(n *element.Node, styleMap map[string]map[string]string) {
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
	height, _ := utils.ConvertToPixels(n.Styles["height"], n.EM, n.Parent.Height)

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

	n.Text.Text = dom.TextContent(n.Node)
	n.Text.Font = f
	n.Text.WordSpacing = int(wordSpacing)
	n.Text.LetterSpacing = int(letterSpacing)

	n.Colors = color.Parse(n.Styles)

	for _, c := range dom.ChildNodes(n.Node) {

		if c.Type == html.ElementNode {
			id := dom.GetAttribute(c, "DOMNODEID")
			node := element.Node{
				Node:   c,
				Parent: n,
				Id:     id,
				Styles: styleMap[id],
			}
			initNodes(&node, styleMap)
			n.Children = append(n.Children, node)
		}
	}
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
	lineHeight, _ := utils.ConvertToPixels(n.Styles["line-height"], n.EM, *width)
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
	n.Text.LineHeight = int(lineHeight)
	n.Text.WhiteSpace = n.Styles["white-space"]
	n.Text.DecorationColor = n.Colors.TextDecoration
	n.Text.DecorationThickness = int(dt)
	n.Text.Overlined = n.Styles["text-decoration"] == "overline"
	n.Text.Underlined = n.Styles["text-decoration"] == "underline"
	n.Text.LineThrough = n.Styles["text-decoration"] == "linethrough"
	n.Text.EM = int(n.EM)

	if n.Styles["word-spacing"] == "" {
		n.Text.WordSpacing = font.MeasureSpace(&n.Text)
	}

	if n.Parent.Width != 0 && n.Styles["display"] != "inline" && n.Styles["width"] == "" {
		*width = (n.Parent.Width - n.Padding.Right) - n.Padding.Left
	} else if n.Styles["width"] == "" {
		lines := font.GetLines(n.Text)
		*width = utils.Max(*width, float32(font.MeasureText(&n.Text, findLongestLine(lines))))

	} else if n.Styles["width"] != "" {
		*width, _ = utils.ConvertToPixels(n.Styles["width"], n.EM, n.Parent.Width)

	}

	n.Text.Width = int(*width)
	*height = font.Render(n)
}

func findLongestLine(lines []string) string {
	var longestLine string
	maxLength := 0

	for _, line := range lines {
		length := len(line)
		if length > maxLength {
			maxLength = length
			longestLine = line
		}
	}

	return longestLine
}
