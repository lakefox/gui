package element

import (
	"fmt"
	"gui/canvas"
	"gui/selector"
	ic "image/color"
	"math"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/font"
)

type Node struct {
	TagName   string
	InnerText string
	InnerHTML string
	OuterHTML string
	Parent    *Node `json:"-"`
	Children  []*Node
	Style     map[string]string
	Id        string
	ClassList ClassList
	Href      string
	Src       string
	Title     string
	Attribute map[string]string
	// !ISSUE: I don't like the scrollx/y
	ScrollX        int
	ScrollY        int
	ScrollLeft     int
	ScrollTop      int
	ScrollHeight   int
	Context        *canvas.Canvas
	PseudoElements map[string]map[string]string

	Value         string
	OnClick       func(Event) `json:"-"`
	OnContextMenu func(Event) `json:"-"`
	OnMouseDown   func(Event) `json:"-"`
	OnMouseUp     func(Event) `json:"-"`
	OnMouseEnter  func(Event) `json:"-"`
	OnMouseLeave  func(Event) `json:"-"`
	OnMouseOver   func(Event) `json:"-"`
	OnMouseMove   func(Event) `json:"-"`
	OnScroll      func(Event) `json:"-"`
	Properties    Properties
}

type State struct {
	// Id         string
	X          float32
	Y          float32
	Z          float32
	Width      float32
	Height     float32
	Border     Border
	Textures   []string
	EM         float32
	Background ic.RGBA
	// Color      ic.RGBA
	Margin       MarginPadding
	Padding      MarginPadding
	Cursor       string
	Crop         Crop
	Hidden       bool
	ScrollHeight int
}

type Crop struct {
	X      int
	Y      int
	Width  int
	Height int
}

// !FLAG: Plan to get rid of this

type Properties struct {
	Id             string
	EventListeners map[string][]func(Event) `json:"-"`
	Events         []string
	Focusable      bool
	Focused        bool
	Editable       bool
	Hover          bool
	Selected       []float32
}

type ClassList struct {
	Classes []string
	Value   string
}

type MarginPadding struct {
	Top    float32
	Left   float32
	Right  float32
	Bottom float32
}

func (c *ClassList) Add(class string) {
	c.Classes = append(c.Classes, class)
	c.Value = strings.Join(c.Classes, " ")
}

func (c *ClassList) Remove(class string) {
	for i, v := range c.Classes {
		if v == class {
			c.Classes = append(c.Classes[:i], c.Classes[i+1:]...)
			break
		}
	}

	c.Value = strings.Join(c.Classes, " ")
}

type Border struct {
	Top    BorderSide
	Right  BorderSide
	Bottom BorderSide
	Left   BorderSide
	Radius BorderRadius
}

type BorderSide struct {
	Width float32
	Style string
	Color ic.RGBA
}

type BorderRadius struct {
	TopLeft     float32
	TopRight    float32
	BottomLeft  float32
	BottomRight float32
}

type Text struct {
	Font                *font.Face
	Color               ic.RGBA
	Text                string
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
	LoadedFont          string
	Last                bool
}

type Shadow struct {
	X     int
	Y     int
	Blur  int
	Color ic.RGBA
}

func (n *Node) GetAttribute(name string) string {
	return n.Attribute[name]
}

func (n *Node) SetAttribute(key, value string) {
	n.Attribute[key] = value
}

func (n *Node) CreateElement(name string) Node {
	return Node{
		TagName:   name,
		InnerText: "",
		OuterHTML: "",
		InnerHTML: "",
		Children:  []*Node{},
		Style:     make(map[string]string),
		Id:        "",
		ClassList: ClassList{
			Classes: []string{},
			Value:   "",
		},
		Href:      "",
		Src:       "",
		Title:     "",
		Attribute: make(map[string]string),
		Value:     "",
		Properties: Properties{
			Id:             "",
			EventListeners: make(map[string][]func(Event)),
			Focusable:      false,
			Focused:        false,
			Editable:       false,
			Hover:          false,
			Selected:       []float32{},
		},
	}
}

func (n *Node) QuerySelectorAll(selectString string) *[]*Node {
	results := []*Node{}
	if TestSelector(selectString, n) {
		results = append(results, n)
	}

	for i := range n.Children {
		el := n.Children[i]
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
		el := n.Children[i]
		cr := el.QuerySelector(selectString)
		if cr.Properties.Id != "" {
			return cr
		}
	}

	return &Node{}
}

func TestSelector(selectString string, n *Node) bool {
	parts := strings.Split(selectString, ">")

	selectors := []string{}
	if n.Properties.Focusable && n.Properties.Focused {
		selectors = append(selectors, ":focus")
	}

	for _, class := range n.ClassList.Classes {
		selectors = append(selectors, "."+class)
	}

	if n.Id != "" {
		selectors = append(selectors, "#"+n.Id)
	}

	selectors = append(selectors, n.TagName)
	part := selector.SplitSelector(strings.TrimSpace(parts[len(parts)-1]))
	// fmt.Println(part)

	has := selector.Contains(part, selectors)

	if len(parts) == 1 || !has {
		return has
	}
	return TestSelector(strings.Join(parts[:len(parts)-1], ">"), n.Parent)
}

var (
	idCounter int64
	mu        sync.Mutex
)

func generateUniqueId(tagName string) string {
	mu.Lock()
	defer mu.Unlock()
	if idCounter == math.MaxInt64 {
		idCounter = 0
	}
	idCounter++
	return tagName + strconv.FormatInt(idCounter, 10)
}

func (n *Node) AppendChild(c *Node) {
	c.Parent = n
	c.Properties.Id = generateUniqueId(c.TagName)
	n.Children = append(n.Children, c)
}

func (n *Node) InsertAfter(c, tgt *Node) {
	c.Parent = n
	c.Properties.Id = generateUniqueId(c.TagName)

	nodeIndex := -1
	for i, v := range n.Children {
		if v.Properties.Id == tgt.Properties.Id {
			nodeIndex = i
			break
		}
	}
	if nodeIndex > -1 {
		n.Children = append(n.Children[:nodeIndex+1], append([]*Node{c}, n.Children[nodeIndex+1:]...)...)
	} else {
		n.AppendChild(c)
	}
}

func (n *Node) InsertBefore(c, tgt *Node) {
	c.Parent = n
	// Set Id

	c.Properties.Id = generateUniqueId(c.TagName)

	nodeIndex := -1
	for i, v := range n.Children {
		if v.Properties.Id == tgt.Properties.Id {
			nodeIndex = i
			break
		}
	}
	if nodeIndex > 0 {
		n.Children = append(n.Children[:nodeIndex], append([]*Node{c}, n.Children[nodeIndex:]...)...)
		// n.Children = append(n.Children, Node{}) // Add a zero value to expand the slice
		// copy(n.Children[nodeIndex+1:], n.Children[nodeIndex:])
		// n.Children[nodeIndex] = c
	} else {
		n.AppendChild(c)
	}

}

func (n *Node) Remove() {
	nodeIndex := -1
	for i, v := range n.Parent.Children {
		if v.Properties.Id == n.Properties.Id {
			nodeIndex = i
			break
		}
	}
	if nodeIndex > 0 {
		n.Parent.Children = append(n.Parent.Children[:nodeIndex], n.Parent.Children[nodeIndex+1:]...)
	}
}

func (n *Node) Focus() {
	if n.Properties.Focusable {
		n.Properties.Focused = true
		n.ClassList.Add(":focus")
	}
}

func (n *Node) Blur() {
	if n.Properties.Focusable {
		n.Properties.Focused = false
		n.ClassList.Remove(":focus")
	}
}

func (n *Node) GetContext(width, height int) *canvas.Canvas {
	n.Style["width"] = strconv.Itoa(width) + "px"
	n.Style["height"] = strconv.Itoa(height) + "px"
	ctx := canvas.NewCanvas(width, height)
	n.Context = ctx
	return ctx
}

func (n *Node) ScrollTo(x, y int) {
	n.ScrollTop = y
	n.ScrollLeft = x
}

type Event struct {
	X           int
	Y           int
	KeyCode     int
	Key         string
	CtrlKey     bool
	MetaKey     bool
	ShiftKey    bool
	AltKey      bool
	Click       bool
	ContextMenu bool
	MouseDown   bool
	MouseUp     bool
	MouseEnter  bool
	MouseLeave  bool
	MouseOver   bool
	KeyUp       bool
	KeyDown     bool
	KeyPress    bool
	Input       bool
	Target      *Node
	Name        string
	Data        interface{}
}

type EventList struct {
	Event Event
	List  []string
}

func (node *Node) AddEventListener(name string, callback func(Event)) {
	if node.Properties.EventListeners == nil {
		node.Properties.EventListeners = make(map[string][]func(Event))
	}
	if node.Properties.EventListeners[name] == nil {
		node.Properties.EventListeners[name] = []func(Event){}
	}
	if !funcInSlice(callback, node.Properties.EventListeners[name]) {
		node.Properties.EventListeners[name] = append(node.Properties.EventListeners[name], callback)
	}
}

func (node *Node) DispatchEvent(event string) {
	// I want to use dispatch event to dispatch a scroll event if the user scrolls too far on a container from crop
	node.Properties.Events = append(node.Properties.Events, event)
}

func funcInSlice(f func(Event), slice []func(Event)) bool {
	for _, item := range slice {
		// Compare function values directly
		if fmt.Sprintf("%p", item) == fmt.Sprintf("%p", f) {
			return true
		}
	}
	return false
}
