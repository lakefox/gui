package cstyle

import (
	"fmt"
	"gui/border"
	"gui/color"
	"gui/element"
	"gui/font"
	"gui/parser"
	"gui/utils"
	"os"
	"sort"
	"strconv"
	"strings"

	imgFont "golang.org/x/image/font"
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
	Handler  func(*element.Node, *CSS) *element.Node
}

type CSS struct {
	Width        float32
	Height       float32
	StyleSheets  []map[string]map[string]string
	Plugins      []Plugin
	Transformers []Transformer
	Document     *element.Node
	Fonts        map[string]imgFont.Face
}

func (c *CSS) Transform(n *element.Node) *element.Node {
	for _, v := range c.Transformers {
		if v.Selector(n) {
			n = v.Handler(n, c)
		}
	}

	for i := 0; i < len(n.Children); i++ {
		tc := c.Transform(n.Children[i])
		// Removing this causes text to break, what does this do??????
		// its because we are using the dom methods to inject on text, not like ulol, prob need to get those working bc they effect user dom as well
		// todo: dom fix, inline-block, text align vert (poss by tgting parent node instead), scroll
		// n = tc.Parent
		n.Children[i] = tc
	}

	return n
}

func (c *CSS) StyleSheet(path string) {
	// Parse the CSS file
	dat, _ := os.ReadFile(path)
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

func (c *CSS) QuickStyles(n *element.Node) map[string]string {
	styles := make(map[string]string)

	// Inherit styles from parent
	if n.Parent != nil {
		ps := n.Parent.Style
		for _, prop := range inheritedProps {
			if value, ok := ps[prop]; ok && value != "" {
				styles[prop] = value
			}
		}
	}

	// Add node's own styles
	for k, v := range n.Style {
		styles[k] = v
	}

	return styles
}

func (c *CSS) GetStyles(n *element.Node) map[string]string {
	styles := make(map[string]string)

	// Inherit styles from parent
	if n.Parent != nil {
		ps := n.Parent.Style
		for _, prop := range inheritedProps {
			if value, ok := ps[prop]; ok && value != "" {
				styles[prop] = value
			}
		}
	}

	// Add node's own styles
	for k, v := range n.Style {
		styles[k] = v
	}

	// Check if node is hovered
	hovered := false
	for _, class := range n.ClassList.Classes {
		if class == ":hover" {
			hovered = true
			break
		}
	}

	// Apply styles from style sheets
	for _, styleSheet := range c.StyleSheets {
		for selector, rules := range styleSheet {
			originalSelector := selector

			if hovered && strings.Contains(selector, ":hover") {
				selector = strings.Replace(selector, ":hover", "", -1)
			}

			if element.TestSelector(selector, n) {
				for k, v := range rules {
					styles[k] = v
				}
			}

			selector = originalSelector // Restore original selector
		}
	}

	// Parse inline styles
	inlineStyles := parser.ParseStyleAttribute(n.GetAttribute("style"))
	for k, v := range inlineStyles {
		styles[k] = v
	}

	// Handle z-index inheritance
	if n.Parent != nil && styles["z-index"] == "" {
		if parentZIndex, ok := n.Parent.Style["z-index"]; ok && parentZIndex != "" {
			z, _ := strconv.Atoi(parentZIndex)
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
	fmt.Printf("Margin: %v\n", self.Margin)
	fmt.Printf("Padding: %v\n", self.Padding)
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

	// Cache the style map
	style := n.Style

	self.Background = color.Parse(style, "background")
	self.Border, _ = border.Parse(style, self, parent)
	border.Draw(&self)

	fs := utils.ConvertToPixels(style["font-size"], parent.EM, parent.Width)
	self.EM = fs

	if style["display"] == "none" {
		self.X, self.Y, self.Width, self.Height = 0, 0, 0, 0
		(*state)[n.Properties.Id] = self
		return n
	}

	// Set Z index value to be sorted in window
	if zIndex, err := strconv.Atoi(style["z-index"]); err == nil {
		self.Z = float32(zIndex)
	}
	if parent.Z > 0 {
		self.Z = parent.Z + 1
	}

	(*state)[n.Properties.Id] = self
	wh := utils.GetWH(*n, state)
	width, height := wh.Width, wh.Height

	x, y := parent.X, parent.Y
	offsetX, offsetY := utils.GetXY(n, state)
	x += offsetX
	y += offsetY

	m := utils.GetMP(*n, wh, state, "margin")
	p := utils.GetMP(*n, wh, state, "padding")
	self.Margin = m
	self.Padding = p

	var top, left, right, bottom bool

	if style["position"] == "absolute" {
		bas := utils.GetPositionOffsetNode(n)
		base := s[bas.Properties.Id]
		if topVal := style["top"]; topVal != "" {
			y = utils.ConvertToPixels(topVal, self.EM, parent.Width) + base.Y
			top = true
		}
		if leftVal := style["left"]; leftVal != "" {
			x = utils.ConvertToPixels(leftVal, self.EM, parent.Width) + base.X
			left = true
		}
		if rightVal := style["right"]; rightVal != "" {
			x = base.Width - width - utils.ConvertToPixels(rightVal, self.EM, parent.Width)
			right = true
		}
		if bottomVal := style["bottom"]; bottomVal != "" {
			y = base.Height - height - utils.ConvertToPixels(bottomVal, self.EM, parent.Width)
			bottom = true
		}
	} else {
		for i, v := range n.Parent.Children {
			if v.Style["position"] != "absolute" {
				if v.Properties.Id == n.Properties.Id {
					if i > 0 {
						sib := n.Parent.Children[i-1]
						sibling := s[sib.Properties.Id]
						if sib.Style["position"] != "absolute" {
							if style["display"] == "inline" {
								y = sibling.Y
								if sib.Style["display"] != "inline" {
									y += sibling.Height
								}
							} else {
								y = sibling.Y + sibling.Height + sibling.Border.Top.Width + sibling.Border.Bottom.Width + sibling.Margin.Bottom
							}
						}
					}
					break
				} else if style["display"] != "inline" {
					vState := s[v.Properties.Id]
					y += vState.Margin.Top + vState.Margin.Bottom + vState.Padding.Top + vState.Padding.Bottom + vState.Height + self.Border.Top.Width
				}
			}
		}
	}

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
		n.InnerText = strings.TrimSpace(n.InnerText)
		self = genTextNode(n, state, c)
	}

	(*state)[n.Properties.Id] = self
	(*state)[n.Parent.Properties.Id] = parent

	var childYOffset float32
	for i := 0; i < len(n.Children); i++ {
		v := n.Children[i]
		v.Parent = n
		n.Children[i] = c.ComputeNodeStyle(v, state)
		cState := (*state)[n.Children[i].Properties.Id]
		if style["height"] == "" && style["min-height"] == "" {
			if v.Style["position"] != "absolute" && cState.Y+cState.Height > childYOffset {
				childYOffset = cState.Y + cState.Height
				self.Height = cState.Y - self.Border.Top.Width - self.Y + cState.Height
				self.Height += cState.Margin.Top + cState.Margin.Bottom + cState.Padding.Top + cState.Padding.Bottom + cState.Border.Top.Width + cState.Border.Bottom.Width
			}
		}
		if cState.Width > self.Width {
			self.Width = cState.Width
		}
	}

	if style["height"] == "" {
		self.Height += self.Padding.Bottom
	}

	(*state)[n.Properties.Id] = self
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Level < plugins[j].Level
	})

	for _, v := range plugins {
		if v.Selector(n) {
			v.Handler(n, state)
		}
	}

	return n
}

func genTextNode(n *element.Node, state *map[string]element.State, css *CSS) element.State {
	s := *state
	self := s[n.Properties.Id]
	parent := s[n.Parent.Properties.Id]

	text := element.Text{}

	bold, italic := false, false
	// !ISSUE: needs bolder and the 100 -> 900
	if n.Style["font-weight"] == "bold" {
		bold = true
	}

	if n.Style["font-style"] == "italic" {
		italic = true
	}

	if text.Font == nil {
		if css.Fonts == nil {
			css.Fonts = map[string]imgFont.Face{}
		}
		fid := n.Style["font-family"] + fmt.Sprint(self.EM, bold, italic)
		if css.Fonts[fid] == nil {
			f, _ := font.LoadFont(n.Style["font-family"], int(self.EM), bold, italic)
			css.Fonts[fid] = f
		}
		fnt := css.Fonts[fid]
		text.Font = &fnt
	}

	letterSpacing := utils.ConvertToPixels(n.Style["letter-spacing"], self.EM, parent.Width)
	wordSpacing := utils.ConvertToPixels(n.Style["word-spacing"], self.EM, parent.Width)
	lineHeight := utils.ConvertToPixels(n.Style["line-height"], self.EM, parent.Width)
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
		dt = utils.ConvertToPixels(n.Style["text-decoration-thickness"], self.EM, parent.Width)
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
	text.Last = n.GetAttribute("last") == "true"

	if n.Style["word-spacing"] == "" {
		text.WordSpacing = font.MeasureSpace(&text)
	}

	img, width := font.Render(&text)
	self.Texture = img

	if n.Style["height"] == "" && n.Style["min-height"] == "" {
		self.Height = float32(text.LineHeight)
	}

	if n.Style["width"] == "" && n.Style["min-width"] == "" {
		self.Width = float32(width)
	}

	return self
}
