package cstyle

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
	"strconv"
	"strings"
)

// !TODO: Make a fine selector to target tags and if it has children or not etc
// + could copy the transformers but idk
type Plugin struct {
	Selector func(*element.Node) bool
	Level    int
	Handler  func(*element.Node, *map[string]element.State)
}

type Transformer struct {
	Selector func(*element.Node) bool
	Handler  func(element.Node, *CSS) element.Node
}

type CSS struct {
	Width        float32
	Height       float32
	StyleSheets  []map[string]map[string]string
	Plugins      []Plugin
	Transformers []Transformer
	Document     *element.Node
}

func (c *CSS) Transform(n element.Node) element.Node {
	for _, v := range c.Transformers {
		if v.Selector(&n) {
			n = v.Handler(n, c)
		}
	}
	for i := 0; i < len(n.Children); i++ {
		v := n.Children[i]
		tc := c.Transform(v)
		n = *tc.Parent
		n.Children[i] = tc
	}

	return n
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
	// "text-align",
	"text-indent",
	"text-justify",
	"text-shadow",
	"text-transform",
	"text-decoration",
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

	// This is different than node.Style
	// temp1 = <span style=​"color:​#a6e22e">​CSS​</span>​
	// temp1.style == CSSStyleDeclaration {0: 'color', accentColor: '', additiveSymbols: '', alignContent: '', alignItems: '', alignSelf: '', …}
	// temp1.getAttribute("style") == 'color:#a6e22e'
	inline := parser.ParseStyleAttribute(n.GetAttribute("style") + ";")
	styles = utils.Merge(styles, inline)
	// add hover and focus css events

	if n.Parent != nil {
		if styles["z-index"] == "" && n.Parent.Style["z-index"] != "" {
			z, _ := strconv.Atoi(n.Parent.Style["z-index"])
			z += 1
			styles["z-index"] = strconv.Itoa(z)
		}
	}

	return styles
}

func (c *CSS) AddPlugin(plugin Plugin) {
	c.Plugins = append(c.Plugins, plugin)
}

func (c *CSS) AddTransformer(transformer Transformer) {
	c.Transformers = append(c.Transformers, transformer)
}

func CheckNode(n *element.Node, state *map[string]element.State) {
	s := *state
	self := s[n.Properties.Id]

	fmt.Println(n.TagName, n.Properties.Id)
	fmt.Printf("ID: %v\n", n.Id)
	fmt.Printf("EM: %v\n", self.EM)
	fmt.Printf("Parent: %v\n", n.Parent.TagName)
	fmt.Printf("Classes: %v\n", n.ClassList.Classes)
	fmt.Printf("Text: %v\n", n.InnerText)
	fmt.Printf("X: %v, Y: %v, Z: %v\n", self.X, self.Y, self.Z)
	fmt.Printf("Width: %v, Height: %v\n", self.Width, self.Height)
	fmt.Printf("Styles: %v\n", n.Style)
	// fmt.Printf("Background: %v\n", self.Background)
	// fmt.Printf("Border: %v\n\n\n", self.Border)
}

func (c *CSS) ComputeNodeStyle(n *element.Node, state *map[string]element.State) *element.Node {

	// Head is not renderable
	if utils.IsParent(*n, "head") {
		return n
	}

	plugins := c.Plugins

	s := *state
	self := s[n.Properties.Id]
	parent := s[n.Parent.Properties.Id]

	self.Background = color.Parse(n.Style, "background")
	self.Border, _ = CompleteBorder(n.Style, self, parent)

	fs, _ := utils.ConvertToPixels(n.Style["font-size"], parent.EM, parent.Width)
	self.EM = fs

	if n.Style["display"] == "none" {
		self.X = 0
		self.Y = 0
		self.Width = 0
		self.Height = 0
		return n
	}

	if n.Style["width"] == "" && n.Style["display"] != "inline" {
		n.Style["width"] = "100%"
	}

	// Set Z index value to be sorted in window
	if n.Style["z-index"] != "" {
		z, _ := strconv.Atoi(n.Style["z-index"])
		self.Z = float32(z)
	}

	if parent.Z > 0 {
		self.Z = parent.Z + 1
	}

	(*state)[n.Properties.Id] = self

	wh := utils.GetWH(*n, state)
	width := wh.Width
	height := wh.Height

	x, y := parent.X, parent.Y
	// !NOTE: Would like to consolidate all XY function into this function like WH
	offsetX, offsetY := utils.GetXY(n, state)
	x += offsetX
	y += offsetY

	var top, left, right, bottom bool = false, false, false, false

	m := utils.GetMP(*n, wh, state, "margin")
	p := utils.GetMP(*n, wh, state, "padding")

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
			if v.Style["position"] != "absolute" {
				if v.Properties.Id == n.Properties.Id {
					if i-1 > 0 {
						sib := n.Parent.Children[i-1]
						sibling := s[sib.Properties.Id]
						if sib.Style["position"] != "absolute" {
							if n.Style["display"] == "inline" {
								if sib.Style["display"] == "inline" {
									y = sibling.Y
								} else {
									y = sibling.Y + sibling.Height
								}
							} else {
								y = sibling.Y + sibling.Height + (sibling.Border.Width * 2) + sibling.Margin.Bottom
							}
						}

					}
					break
				} else if n.Style["display"] != "inline" {
					vState := s[v.Properties.Id]
					y += vState.Margin.Top + vState.Margin.Bottom + vState.Padding.Top + vState.Padding.Bottom + vState.Height + (self.Border.Width)
				}
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

	self.X = x
	self.Y = y
	self.Width = width
	self.Height = height
	(*state)[n.Properties.Id] = self

	if !utils.ChildrenHaveText(n) && len(n.InnerText) > 0 {
		// Confirm text exists
		if len(strings.TrimSpace(n.InnerText)) > 0 {
			n.InnerText = strings.TrimSpace(n.InnerText)
			self = genTextNode(n, state)
		}
	}

	(*state)[n.Properties.Id] = self
	(*state)[n.Parent.Properties.Id] = parent

	// Call children here
	// !TODO: Make something that stops rendering things out of fov

	if self.Y < c.Height {
		var childYOffset float32
		for i := 0; i < len(n.Children); i++ {
			v := n.Children[i]
			v.Parent = n
			// This is were the tainting comes from
			n.Children[i] = *c.ComputeNodeStyle(&v, state)

			cState := (*state)[n.Children[i].Properties.Id]
			if n.Style["height"] == "" {
				if v.Style["position"] != "absolute" && cState.Y+cState.Height > childYOffset {
					childYOffset = cState.Y + cState.Height
					self.Height = (cState.Y - self.Border.Width) - (self.Y) + cState.Height
					self.Height += cState.Margin.Top
					self.Height += cState.Margin.Bottom
					self.Height += cState.Padding.Top
					self.Height += cState.Padding.Bottom
				}
			}
			if cState.Width > self.Width {
				self.Width = cState.Width
			}
		}
	} else {
		return n
	}

	self.Height += self.Padding.Bottom

	(*state)[n.Properties.Id] = self

	// Sorting the array by the Level field
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Level < plugins[j].Level
	})

	for _, v := range plugins {
		if v.Selector(n) {
			v.Handler(n, state)
		}
	}

	// CheckNode(n, state)

	return n
}

func CompleteBorder(cssProperties map[string]string, self, parent element.State) (element.Border, error) {
	// Split the shorthand into components
	borderComponents := strings.Fields(cssProperties["border"])

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

		w, _ := utils.ConvertToPixels(width, self.EM, parent.Width)

		return element.Border{
			Width:  w,
			Style:  style,
			Color:  parsedColor,
			Radius: cssProperties["border-radius"],
		}, nil
	}

	return element.Border{}, fmt.Errorf("invalid border shorthand format")
}

func genTextNode(n *element.Node, state *map[string]element.State) element.State {
	s := *state
	self := s[n.Properties.Id]
	parent := s[n.Parent.Properties.Id]

	text := element.Text{}

	bold, italic := false, false

	if n.Style["font-weight"] == "bold" {
		bold = true
	}

	if n.Style["font-style"] == "italic" {
		italic = true
	}

	if text.Font == nil {
		f, _ := font.LoadFont(n.Style["font-family"], int(self.EM), bold, italic)
		text.Font = f
	}

	letterSpacing, _ := utils.ConvertToPixels(n.Style["letter-spacing"], self.EM, parent.Width)
	wordSpacing, _ := utils.ConvertToPixels(n.Style["word-spacing"], self.EM, parent.Width)
	lineHeight, _ := utils.ConvertToPixels(n.Style["line-height"], self.EM, parent.Width)
	if lineHeight == 0 {
		lineHeight = self.EM + 3
	}

	text.LineHeight = int(lineHeight)
	text.WordSpacing = int(wordSpacing)
	text.LetterSpacing = int(letterSpacing)
	wb := " "

	if n.Style["word-wrap"] == "break-word" {
		wb = ""
	}

	if n.Style["text-wrap"] == "wrap" || n.Style["text-wrap"] == "balance" {
		wb = ""
	}

	var dt float32

	if n.Style["text-decoration-thickness"] == "auto" || n.Style["text-decoration-thickness"] == "" {
		dt = self.EM / 7
	} else {
		dt, _ = utils.ConvertToPixels(n.Style["text-decoration-thickness"], self.EM, parent.Width)
	}

	col := color.Parse(n.Style, "font")

	self.Color = col

	text.Color = col
	text.DecorationColor = color.Parse(n.Style, "decoration")
	text.Align = n.Style["text-align"]
	text.WordBreak = wb
	text.WordSpacing = int(wordSpacing)
	text.LetterSpacing = int(letterSpacing)
	text.WhiteSpace = n.Style["white-space"]
	text.DecorationThickness = int(dt)
	text.Overlined = n.Style["text-decoration"] == "overline"
	text.Underlined = n.Style["text-decoration"] == "underline"
	text.LineThrough = n.Style["text-decoration"] == "linethrough"
	text.EM = int(self.EM)
	text.Width = int(parent.Width)
	text.Text = n.InnerText

	if n.Style["word-spacing"] == "" {
		text.WordSpacing = font.MeasureSpace(&text)
	}

	img, width := font.Render(&text)
	self.Texture = img

	if n.Style["height"] == "" {
		self.Height = float32(text.LineHeight)
	}

	if n.Style["width"] == "" {
		self.Width = float32(width)
	}

	return self
}
