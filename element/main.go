package element

import (
	"image"
	ic "image/color"

	"golang.org/x/image/font"

	"golang.org/x/net/html"
)

type Node struct {
	Node     *html.Node
	Parent   *Node
	Children []Node
	Styles   map[string]string
	Id       string
	X        float32
	Y        float32
	Width    float32
	Height   float32
	Margin   Margin
	Padding  Padding
	Border   Border
	EM       float32
	Text     Text
	Colors   Colors
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
