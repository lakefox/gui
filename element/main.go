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
	Node        *html.Node
	Type        html.NodeType
	TagName     string
	Parent      *Node
	Children    []Node
	Styles      map[string]string
	Id          string
	X           float32
	Y           float32
	Width       float32
	Height      float32
	Margin      Margin
	Padding     Padding
	Border      Border
	EM          float32
	Text        Text
	Colors      Colors
	PrevSibling *Node
	NextSibling *Node
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
	Text                string
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

	for _, attr := range n.Node.Attr {
		attributes[attr.Key] = attr.Val
	}
	return attributes[name]
}

func (n *Node) SetAttribute(key, value string) {
	// Iterate through the attributes
	for i, attr := range n.Node.Attr {
		// If the attribute key matches, update its value
		if attr.Key == key {
			n.Node.Attr[i].Val = value
			return
		}
	}

	// If the attribute key was not found, add a new attribute
	n.Node.Attr = append(n.Node.Attr, html.Attribute{
		Key: key,
		Val: value,
	})
}

func (n *Node) QuerySelectorAll(selectString string) []Node {
	results := []Node{}
	if TestSelector(selectString, n) {
		results = append(results, *n)
	}

	for _, v := range n.Children {
		cr := v.QuerySelectorAll(selectString)
		if len(cr) > 0 {
			results = append(results, cr...)
		}
	}
	return results
}

func (n *Node) QuerySelector(selectString string) Node {
	if TestSelector(selectString, n) {
		return *n
	}

	for _, v := range n.Children {
		cr := v.QuerySelector(selectString)
		if cr.Id != "" {
			return cr
		}
	}

	return Node{}
}

func (node *Node) InnerText() string {
	var result strings.Builder

	var getText func(*html.Node)
	getText = func(n *html.Node) {
		if n.Type == html.TextNode {
			result.WriteString(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getText(c)
		}
	}

	getText(node.Node)

	return result.String()
}

func TestSelector(selectString string, n *Node) bool {
	parts := strings.Split(selectString, ">")

	selectors := selector.GetCSSSelectors(n.Node, []string{})

	part := selector.SplitSelector(strings.TrimSpace(parts[len(parts)-1]))

	has := selector.Contains(part, selectors)

	if len(parts) == 1 || !has {
		return has
	} else {
		return TestSelector(strings.Join(parts[0:len(parts)-1], ">"), n.Parent)
	}
}
