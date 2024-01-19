package element

import (
	"gui/selector"
	"image"
	ic "image/color"
	"strings"

	"golang.org/x/image/font"

	"golang.org/x/net/html"
)

type Node struct {
	TagName       string
	InnerText     string
	Parent        *Node
	Children      []Node
	Style         map[string]string
	PrevSibling   *Node
	NextSibling   *Node
	OnClick       func(Event)
	OnContextMenu func(Event)
	OnMouseDown   func(Event)
	OnMouseUp     func(Event)
	OnMouseEnter  func(Event)
	OnMouseLeave  func(Event)
	OnMouseOver   func(Event)
	OnMouseMove   func(Event)
	Properties    Properties
}

type Properties struct {
	Node           *html.Node
	Type           html.NodeType
	Id             string
	X              float32
	Y              float32
	Width          float32
	Height         float32
	Margin         Margin
	Padding        Padding
	Border         Border
	EventListeners map[string][]func(Event)
	EM             float32
	Text           Text
	Colors         Colors
}

type Margin struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

type Padding struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

type Border struct {
	Width  string
	Style  string
	Color  ic.RGBA
	Radius string
}

type Text struct {
	Font                font.Face
	Color               ic.RGBA
	Image               *image.RGBA
	Underlined          bool
	Overlined           bool
	LineThrough         bool
	DecorationColor     ic.RGBA
	DecorationThickness int
	Align               string
	Indent              int // very low priority
	LetterSpacing       int
	LineHeight          int
	WordSpacing         int
	WhiteSpace          string
	Shadows             []Shadow // need
	Width               int
	WordBreak           string
	EM                  int
	X                   int
}

type Shadow struct {
	X     int
	Y     int
	Blur  int
	Color ic.RGBA
}

// Color represents an RGBA color
type Colors struct {
	Background     ic.RGBA
	Font           ic.RGBA
	TextDecoration ic.RGBA
}

func (n *Node) GetAttribute(name string) string {
	attributes := make(map[string]string)

	for _, attr := range n.Properties.Node.Attr {
		attributes[attr.Key] = attr.Val
	}
	return attributes[name]
}

func (n *Node) SetAttribute(key, value string) {
	// Iterate through the attributes
	for i, attr := range n.Properties.Node.Attr {
		// If the attribute key matches, update its value
		if attr.Key == key {
			n.Properties.Node.Attr[i].Val = value
			return
		}
	}

	// If the attribute key was not found, add a new attribute
	n.Properties.Node.Attr = append(n.Properties.Node.Attr, html.Attribute{
		Key: key,
		Val: value,
	})
}

func (n *Node) QuerySelectorAll(selectString string) *[]*Node {
	results := []*Node{}
	if TestSelector(selectString, n) {
		results = append(results, n)
	}

	for i := range n.Children {
		el := &n.Children[i]
		cr := el.QuerySelectorAll(selectString)
		if len(*cr) > 0 {
			results = append(results, *cr...)
		}
	}
	return &results
}

func (n *Node) QuerySelector(selectString string) *Node {
	if TestSelector(selectString, n) {
		return n
	}

	for i := range n.Children {
		el := &n.Children[i]
		cr := el.QuerySelector(selectString)
		if cr.Properties.Id != "" {
			return cr
		}
	}

	return &Node{}
}

func TestSelector(selectString string, n *Node) bool {
	parts := strings.Split(selectString, ">")

	selectors := selector.GetCSSSelectors(n.Properties.Node, []string{})

	part := selector.SplitSelector(strings.TrimSpace(parts[len(parts)-1]))

	has := selector.Contains(part, selectors)

	if len(parts) == 1 || !has {
		return has
	} else {
		return TestSelector(strings.Join(parts[0:len(parts)-1], ">"), n.Parent)
	}
}

type Event struct {
	X           int
	Y           int
	Click       bool
	ContextMenu bool
	MouseDown   bool
	MouseUp     bool
	MouseEnter  bool
	MouseLeave  bool
	MouseOver   bool
	Target      Node
}

func (node *Node) AddEventListener(name string, callback func(Event)) {
	if node.Properties.EventListeners == nil {
		node.Properties.EventListeners = make(map[string][]func(Event))
	}
	if node.Properties.EventListeners[name] == nil {
		node.Properties.EventListeners[name] = []func(Event){}
	}
	node.Properties.EventListeners[name] = append(node.Properties.EventListeners[name], callback)
}
