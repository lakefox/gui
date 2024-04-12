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
	"os"
	"slices"
	"sort"
	"strings"
)

type Plugin struct {
	Styles  map[string]string
	Level   int
	Handler func(*element.Node, *map[string]element.State)
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

func (c *CSS) GetStyles(n element.Node) map[string]string {
	styles := map[string]string{}

	if n.Parent != nil {
		ps := c.GetStyles(*n.Parent)
		for _, v := range inheritedProps {
			if ps[v] != "" {
				styles[v] = ps[v]
			}
		}
	}
	for k, v := range n.Style {
		styles[k] = v
	}
	hovered := false
	if slices.Contains(n.ClassList.Classes, ":hover") {
		hovered = true
	}

	for _, styleSheet := range c.StyleSheets {
		for selector := range styleSheet {
			// fmt.Println(selector, n.Properties.Id)
			key := selector
			if strings.Contains(selector, ":hover") && hovered {
				selector = strings.Replace(selector, ":hover", "", -1)
			}
			if element.TestSelector(selector, &n) {
				for k, v := range styleSheet[key] {
					styles[k] = v
				}
			}

		}
	}

	// !FLAG: why is this needed, the "attribute" is n.Style that should be mapped during init
	inline := parser.ParseStyleAttribute(n.GetAttribute("style") + ";")
	styles = utils.Merge(styles, inline)
	// add hover and focus css events

	return styles
}

func (c *CSS) AddPlugin(plugin Plugin) {
	c.Plugins = append(c.Plugins, plugin)
}

func CheckNode(n *element.Node, state *map[string]element.State) {
	s := *state
	self := s[n.Properties.Id]

	fmt.Println(n.TagName, n.Properties.Id)
	fmt.Printf("ID: %v\n", n.Id)
	fmt.Printf("Parent: %v\n", n.Parent.TagName)
	fmt.Printf("Classes: %v\n", n.ClassList.Classes)
	fmt.Printf("Text: %v\n", self.Text.Text)
	fmt.Printf("X: %v, Y: %v\n", self.X, self.Y)
	fmt.Printf("Width: %v, Height: %v\n", self.Width, self.Height)
	fmt.Printf("Styles: %v\n", n.Style)
	fmt.Printf("Background: %v\n", self.Background)
	fmt.Printf("Border: %v\n\n\n", self.Border)
}

func (c *CSS) ComputeNodeStyle(n *element.Node, state *map[string]element.State) *element.Node {
	// Head is not renderable
	if utils.IsParent(*n, "head") {
		return n
	}
	plugins := c.Plugins
	// !FLAG: This should add to state.Style instead as the element.Node should be un effected by the engine
	// 		  currently this adds styles to the style attribute that the use did not explisitly set

	n.Style = c.GetStyles(*n)
	s := *state
	self := s[n.Properties.Id]
	parent := s[n.Parent.Properties.Id]

	self.Background = color.Parse(n.Style, "background")
	self.Border, _ = CompleteBorder(n.Style)

	fs, _ := utils.ConvertToPixels(n.Style["font-size"], parent.EM, parent.Width)
	self.EM = fs

	if n.Style["display"] == "none" {
		self.X = 0
		self.Y = 0
		self.Width = 0
		self.Height = 0
		return n
	}

	wh := utils.GetWH(*n)
	width := wh.Width
	height := wh.Height

	x, y := parent.X, parent.Y

	var top, left, right, bottom bool = false, false, false, false

	m := utils.GetMP(*n, "margin")
	p := utils.GetMP(*n, "padding")

	self.Margin = m
	self.Padding = p

	if n.Style["position"] == "absolute" {
		bas := utils.GetPositionOffsetNode(n)
		base := s[bas.Properties.Id]
		if n.Style["top"] != "" {
			v, _ := utils.ConvertToPixels(n.Style["top"], self.EM, parent.Width)
			y = v + base.Y
			top = true
		}
		if n.Style["left"] != "" {
			v, _ := utils.ConvertToPixels(n.Style["left"], self.EM, parent.Width)
			x = v + base.X
			left = true
		}
		if n.Style["right"] != "" {
			v, _ := utils.ConvertToPixels(n.Style["right"], self.EM, parent.Width)
			x = (base.Width - width) - v
			right = true
		}
		if n.Style["bottom"] != "" {
			v, _ := utils.ConvertToPixels(n.Style["bottom"], self.EM, parent.Width)
			y = (base.Height - height) - v
			bottom = true
		}
	} else {
		for i, v := range n.Parent.Children {
			if v.Properties.Id == n.Properties.Id {
				if i-1 > 0 {
					sib := n.Parent.Children[i-1]
					sibling := s[sib.Properties.Id]
					if n.Style["display"] == "inline" {
						if sib.Style["display"] == "inline" {
							y = sibling.Y
						} else {
							y = sibling.Y + sibling.Height
						}
					} else {
						y = sibling.Y + sibling.Height
					}
				}
				break
			} else if n.Style["display"] != "inline" {
				vState := s[v.Properties.Id]
				y += vState.Margin.Top + vState.Margin.Bottom + vState.Padding.Top + vState.Padding.Bottom + vState.Height
			}
		}
	}

	// Display modes need to be calculated here

	relPos := !top && !left && !right && !bottom

	if left || relPos {
		x += m.Left
	}
	if top || relPos {
		y += m.Top
	}
	if right {
		x -= m.Right
	}
	if bottom {
		y -= m.Bottom
	}

	// fmt.Println(n.InnerText, len(n.Children))

	if !utils.ChildrenHaveText(n) {
		// Confirm text exists
		if len(n.InnerText) > 0 {
			innerWidth := width
			innerHeight := height
			(*state)[n.Properties.Id] = self
			self = genTextNode(n, &innerWidth, &innerHeight, p, state)
			width = innerWidth + p.Left + p.Right
			height = innerHeight
		}
	}

	self.X = x
	self.Y = y
	self.Width = width
	self.Height = height

	(*state)[n.Properties.Id] = self
	(*state)[n.Parent.Properties.Id] = parent

	// CheckNode(n, state)

	// Call children here

	var childYOffset float32
	for i, v := range n.Children {
		v.Parent = n
		n.Children[i] = *c.ComputeNodeStyle(&v, state)
		if n.Style["height"] == "" {
			cState := s[n.Children[i].Properties.Id]
			if n.Children[i].Style["position"] != "absolute" && cState.Y > childYOffset {
				childYOffset = cState.Y
				self.Height += cState.Height
				self.Height += cState.Margin.Top
				self.Height += cState.Margin.Bottom
				self.Height += cState.Padding.Top
				self.Height += cState.Padding.Bottom
			}

		}
	}

	(*state)[n.Properties.Id] = self

	// Sorting the array by the Level field
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Level < plugins[j].Level
	})

	for _, v := range plugins {
		matches := true
		for name, value := range v.Styles {
			if n.Style[name] != value && !(value == "*") {
				matches = false
			}
		}
		if matches {
			v.Handler(n, state)
		}
	}

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

func genTextNode(n *element.Node, width, height *float32, p element.MarginPadding, state *map[string]element.State) element.State {
	s := *state
	self := s[n.Properties.Id]
	parent := s[n.Parent.Properties.Id]

	bold, italic := false, false

	if n.Style["font-weight"] == "bold" {
		bold = true
	}

	if n.Style["font-style"] == "italic" {
		italic = true
	}

	if self.Text.Font == nil {
		f, _ := font.LoadFont(n.Style["font-family"], int(self.EM), bold, italic)
		self.Text.Font = f
	}

	letterSpacing, _ := utils.ConvertToPixels(n.Style["letter-spacing"], self.EM, *width)
	wordSpacing, _ := utils.ConvertToPixels(n.Style["word-spacing"], self.EM, *width)
	lineHeight, _ := utils.ConvertToPixels(n.Style["line-height"], self.EM, *width)
	if lineHeight == 0 {
		lineHeight = self.EM + 3
	}

	self.Text.LineHeight = int(lineHeight)
	self.Text.WordSpacing = int(wordSpacing)
	self.Text.LetterSpacing = int(letterSpacing)
	wb := " "

	if n.Style["word-wrap"] == "break-word" {
		wb = ""
	}

	if n.Style["text-wrap"] == "wrap" || n.Style["text-wrap"] == "balance" {
		wb = ""
	}

	var dt float32

	if n.Style["text-decoration-thickness"] == "auto" || n.Style["text-decoration-thickness"] == "" {
		dt = 3
	} else {
		dt, _ = utils.ConvertToPixels(n.Style["text-decoration-thickness"], self.EM, *width)
	}

	col := color.Parse(n.Style, "font")

	self.Text.Color = col
	self.Text.DecorationColor = color.Parse(n.Style, "decoration")
	self.Text.Align = n.Style["text-align"]
	self.Text.WordBreak = wb
	self.Text.WordSpacing = int(wordSpacing)
	self.Text.LetterSpacing = int(letterSpacing)
	self.Text.WhiteSpace = n.Style["white-space"]
	self.Text.DecorationThickness = int(dt)
	self.Text.Overlined = n.Style["text-decoration"] == "overline"
	self.Text.Underlined = n.Style["text-decoration"] == "underline"
	self.Text.LineThrough = n.Style["text-decoration"] == "linethrough"
	self.Text.EM = int(self.EM)
	self.Text.Width = int(parent.Width)
	self.Text.Text = n.InnerText

	if n.Style["word-spacing"] == "" {
		self.Text.WordSpacing = font.MeasureSpace(&self.Text)
	}
	if parent.Width != 0 && n.Style["display"] != "inline" && n.Style["width"] == "" {
		*width = (parent.Width - p.Right) - p.Left
	} else if n.Style["width"] == "" {
		*width = utils.Max(*width, float32(font.MeasureLongest(&self)))
	} else if n.Style["width"] != "" {
		*width, _ = utils.ConvertToPixels(n.Style["width"], self.EM, parent.Width)
	}

	self.Text.Width = int(*width)
	self.Width = *width
	fmt.Println(n.TagName, n.Style["width"], *width)
	h := font.Render(&self)
	if n.Style["height"] == "" {
		*height = h
	}

	return self

}
