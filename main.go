package gui

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	adapter "gui/adapters"
	"gui/cstyle"
	"gui/cstyle/plugins/flex"
	"gui/cstyle/plugins/inline"
	"gui/cstyle/plugins/textAlign"
	flexprep "gui/cstyle/transformers/flex"
	"gui/cstyle/transformers/ol"
	"gui/cstyle/transformers/text"
	"gui/cstyle/transformers/ul"
	"gui/font"
	"gui/library"
	"gui/scripts"
	"gui/scripts/a"
	"image"
	"time"

	"gui/element"
	"gui/events"
	"gui/utils"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	imgFont "golang.org/x/image/font"

	"golang.org/x/net/html"
)

// _ "net/http/pprof"

//go:embed master.css
var mastercss string

type Window struct {
	CSS      cstyle.CSS
	Document element.Node
	Adapter  *adapter.Adapter
	Scripts  scripts.Scripts
}

func Open(path string) Window {
	window := New()

	styleSheets, styleTags, htmlNodes := parseHTMLFromFile(path)

	for _, v := range styleSheets {
		window.CSS.StyleSheet(v)
	}

	for _, v := range styleTags {
		window.CSS.StyleTag(v)
	}

	CreateNode(htmlNodes, &window.Document)

	return window
}

func New() Window {
	css := cstyle.CSS{
		Width:  800,
		Height: 450,
	}

	css.StyleTag(mastercss)
	// This is still apart of computestyle
	// css.AddPlugin(position.Init())
	css.AddPlugin(inline.Init())
	// css.AddPlugin(block.Init())
	css.AddPlugin(textAlign.Init())
	// css.AddPlugin(inlineText.Init())
	css.AddPlugin(flex.Init())

	css.AddTransformer(flexprep.Init())
	css.AddTransformer(ul.Init())
	css.AddTransformer(ol.Init())
	// css.AddTransformer(textInline.Init())
	css.AddTransformer(text.Init())

	el := element.Node{}
	document := el.CreateElement("ROOT")
	document.Style["width"] = "800px"
	document.Style["height"] = "450px"
	document.Properties.Id = "ROOT"

	s := scripts.Scripts{}
	s.Add(a.Init())

	return Window{
		CSS:      css,
		Document: document,
		Scripts:  s,
	}
}

func (w *Window) Render(doc *element.Node, state *map[string]element.State) []element.State {
	s := *state

	flatDoc := flatten(doc)

	store := []element.State{}

	keys := []string{}

	for _, v := range flatDoc {
		store = append(store, s[v.Properties.Id])
		keys = append(keys, v.Properties.Id)
	}

	// Create a set of keys to keep
	keysSet := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		keysSet[key] = struct{}{}
	}

	// Iterate over the map and delete keys not in the set
	for k := range s {
		if _, found := keysSet[k]; !found {
			delete(s, k)
		}
	}

	return store
}

func flatten(n *element.Node) []*element.Node {
	var nodes []*element.Node
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

func View(data *Window, width, height int) {

	shelf := library.Shelf{
		Textures:   map[string]*image.RGBA{},
		References: map[string]bool{},
	}

	debug := false
	data.Document.Style["width"] = strconv.Itoa(int(width)) + "px"
	data.Document.Style["height"] = strconv.Itoa(int(height)) + "px"

	data.Adapter.Library = &shelf
	data.Adapter.Init(width, height)
	// wm.FPSCounterOn = true

	evts := map[string]element.EventList{}

	eventStore := &evts

	state := map[string]element.State{}
	state["ROOT"] = element.State{
		Width:  float32(width),
		Height: float32(height),
	}

	shouldStop := false

	var hash []byte
	var rd []element.State

	// Load init font
	if data.CSS.Fonts == nil {
		data.CSS.Fonts = map[string]imgFont.Face{}
	}
	fid := "Georgia 16px false false"
	if data.CSS.Fonts[fid] == nil {
		f, _ := font.LoadFont("Georgia", 16, false, false)
		data.CSS.Fonts[fid] = f
	}

	newWidth, newHeight := width, height

	data.Adapter.AddEventListener("windowresize", func(e element.Event) {
		wh := e.Data.(map[string]int)
		newWidth = wh["width"]
		newHeight = wh["height"]
	})

	data.Adapter.AddEventListener("close", func(e element.Event) {
		shouldStop = true
	})

	data.Adapter.AddEventListener("keydown", func(e element.Event) {
		fmt.Println("Down: ", e.Data.(int))
	})
	data.Adapter.AddEventListener("keyup", func(e element.Event) {
		fmt.Println("Up: ", e.Data.(int))
	})

	// data.Adapter.AddEventListener("mousemove", func(e element.Event) {
	// 	pos := e.Data.([]int)
	// 	fmt.Println("Mouse: ", pos)
	// })

	data.Adapter.AddEventListener("scroll", func(e element.Event) {
		fmt.Println("Scroll: ", e.Data.(int))
	})

	// Main game loop
	for !shouldStop {

		if !shouldStop && debug {
			shouldStop = true
		}
		// Check if the window size has changed

		resize := false

		if newWidth != width || newHeight != height {
			resize = true
			// Window has been resized, handle the event
			width = newWidth
			height = newHeight

			data.CSS.Width = float32(width)
			data.CSS.Height = float32(height)

			data.Document.Style["width"] = strconv.Itoa(int(width)) + "px"
			data.Document.Style["height"] = strconv.Itoa(int(height)) + "px"
		}

		newHash, _ := hashStruct(&data.Document.Children[0])
		eventStore = events.GetEvents(data.Document.Children[0], &state, eventStore)
		if !bytes.Equal(hash, newHash) || resize {

			hash = newHash
			lastChange := time.Now()
			fmt.Println("########################")
			lastChange1 := time.Now()

			// newDoc := data.Document.Children[0]                                     // speed up
			// change to add styles
			newDoc := CopyNode(data.CSS, data.Document.Children[0], &data.Document) // speed up
			fmt.Println("Copy Node: ", time.Since(lastChange1))
			lastChange1 = time.Now()
			newDoc = data.CSS.Transform(newDoc)
			fmt.Println("Transform: ", time.Since(lastChange1))
			lastChange1 = time.Now()

			state["ROOT"] = element.State{
				Width:  float32(width),
				Height: float32(height),
			}

			data.CSS.ComputeNodeStyle(newDoc, &state, &shelf) // speed up
			fmt.Println("Compute Node Style: ", time.Since(lastChange1))
			lastChange1 = time.Now()

			rd = data.Render(newDoc, &state)
			fmt.Println("Render: ", time.Since(lastChange1))
			lastChange1 = time.Now()

			data.Adapter.Load(rd) // speed up
			fmt.Println("Load: ", time.Since(lastChange1))
			lastChange1 = time.Now()

			AddHTML(&data.Document)
			fmt.Println("Add HTML: ", time.Since(lastChange1))

			data.Scripts.Run(&data.Document)
			fmt.Println("#", time.Since(lastChange))
			shelf.Close()
		}
		data.Adapter.Render(rd)

		// could use a return value that indicates whether or not a event has ran to ramp/deramp fps based on activity

		events.RunEvents(eventStore)
	}
}

func CopyNode(c cstyle.CSS, node *element.Node, parent *element.Node) *element.Node {
	// n := element.Node{
	// 	TagName:   node.TagName,
	// 	InnerText: node.InnerText,
	// 	Style:     node.Style,
	// 	Id:        node.Id,
	// 	ClassList: node.ClassList,
	// 	Href:      node.Href,
	// 	Src:       node.Src,
	// 	Title:     node.Title,
	// 	Attribute: node.Attribute,
	// 	Value:     node.Value,
	// 	ScrollY:   node.ScrollY,
	// 	InnerHTML: node.InnerHTML,
	// 	OuterHTML: node.OuterHTML,
	// 	Parent:    parent,
	// 	Properties: element.Properties{
	// 		Id:        node.Properties.Id,
	// 		Focusable: node.Properties.Focusable,
	// 		Focused:   node.Properties.Focused,
	// 		Editable:  node.Properties.Editable,
	// 		Hover:     node.Properties.Hover,
	// 		Selected:  node.Properties.Selected,
	// 	},
	// }
	n := *node
	n.Parent = parent
	// lastChange1 := time.Now()
	// for get styles, pre load all of the selectors into a map or slice then use a map to point to the index in the slice then search that way
	n.Style = c.GetStyles(&n)
	// fmt.Println("Styles: ", time.Since(lastChange1))

	if len(node.Children) > 0 {
		n.Children = make([]*element.Node, 0, len(node.Children))
		for _, v := range node.Children {
			n.Children = append(n.Children, CopyNode(c, v, &n))
		}
	}

	return &n
}

func CreateNode(node *html.Node, parent *element.Node) {
	if node.Type == html.ElementNode {
		newNode := parent.CreateElement(node.Data)
		for _, attr := range node.Attr {
			if attr.Key == "class" {
				classes := strings.Split(attr.Val, " ")
				for _, class := range classes {
					newNode.ClassList.Add(class)
				}
			} else if attr.Key == "id" {
				newNode.Id = attr.Val
			} else if attr.Key == "contenteditable" && (attr.Val == "" || attr.Val == "true") {
				newNode.Properties.Editable = true
			} else if attr.Key == "href" {
				newNode.Href = attr.Val
			} else if attr.Key == "src" {
				newNode.Src = attr.Val
			} else if attr.Key == "title" {
				newNode.Title = attr.Val
			} else {
				newNode.SetAttribute(attr.Key, attr.Val)
			}
		}
		newNode.InnerText = strings.TrimSpace(utils.GetInnerText(node))
		// Recursively traverse child nodes
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode {
				CreateNode(child, &newNode)
			}
		}
		parent.AppendChild(&newNode)

	} else {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode {
				CreateNode(child, parent)
			}
		}
	}
}

func AddHTML(n *element.Node) {
	// Head is not renderable
	n.InnerHTML = utils.InnerHTML(n)
	tag, closing := utils.NodeToHTML(n)
	n.OuterHTML = tag + n.InnerHTML + closing
	for i := range n.Children {
		AddHTML(n.Children[i])
	}
}

func parseHTMLFromFile(path string) ([]string, []string, *html.Node) {
	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var htmlContent string

	for scanner.Scan() {
		htmlContent += scanner.Text() + "\n"
	}

	htmlContent = removeHTMLComments(htmlContent)

	doc, _ := html.Parse(strings.NewReader(encapsulateText(removeWhitespaceBetweenTags(htmlContent))))

	// Extract stylesheet link tags and style tags
	stylesheets := extractStylesheets(doc, filepath.Dir(path))
	styleTags := extractStyleTags(doc)

	return stylesheets, styleTags, doc
}

func extractStylesheets(n *html.Node, baseDir string) []string {
	var stylesheets []string

	var dfs func(*html.Node)
	dfs = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "link" {
			var href string
			isStylesheet := false

			for _, attr := range node.Attr {
				if attr.Key == "rel" && attr.Val == "stylesheet" {
					isStylesheet = true
				} else if attr.Key == "href" {
					href = attr.Val
				}
			}

			if isStylesheet {
				resolvedHref := localizePath(baseDir, href)
				stylesheets = append(stylesheets, resolvedHref)
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}

	dfs(n)
	return stylesheets
}

func extractStyleTags(n *html.Node) []string {
	var styleTags []string

	var dfs func(*html.Node)
	dfs = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "style" {
			var styleContent strings.Builder
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					styleContent.WriteString(c.Data)
				}
			}
			styleTags = append(styleTags, styleContent.String())
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}

	dfs(n)
	return styleTags
}

func localizePath(rootPath, filePath string) string {
	// Check if the file path has a scheme, indicating it's a URL
	u, err := url.Parse(filePath)
	if err == nil && u.Scheme != "" {
		return filePath
	}

	// Join the root path and the file path to create an absolute path
	absPath := filepath.Join(rootPath, filePath)

	// If the absolute path is the same as the original path, return it
	if absPath == filePath {
		return filePath
	}

	return "./" + absPath
}

func encapsulateText(htmlString string) string {
	openOpen := regexp.MustCompile(`(<\w+[^>]*>)([^<]+)(<\w+[^>]*>)`)
	closeOpen := regexp.MustCompile(`(</\w+[^>]*>)([^<]+)(<\w+[^>]*>)`)
	closeClose := regexp.MustCompile(`(<\/\w+[^>]*>)([^<]+)(<\/\w+[^>]*>)`)
	a := matchFactory(openOpen)
	t := openOpen.ReplaceAllStringFunc(htmlString, a)
	// fmt.Println(t)
	b := matchFactory(closeOpen)
	u := closeOpen.ReplaceAllStringFunc(t, b)
	// fmt.Println(u)
	c := matchFactory(closeClose)
	v := closeClose.ReplaceAllStringFunc(u, c)
	// fmt.Println(v)
	return v
}

func matchFactory(re *regexp.Regexp) func(string) string {
	return func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) != 4 {
			return match
		}

		// Process submatches
		if len(removeWhitespace(submatches[2])) > 0 {
			return submatches[1] + "<notaspan>" + submatches[2] + "</notaspan>" + submatches[3]
		} else {
			return match
		}
	}
}
func removeWhitespace(htmlString string) string {
	// Remove extra white space
	reSpaces := regexp.MustCompile(`\s+`)
	htmlString = reSpaces.ReplaceAllString(htmlString, " ")

	// Trim leading and trailing white space
	htmlString = strings.TrimSpace(htmlString)

	return htmlString
}

func removeHTMLComments(htmlString string) string {
	re := regexp.MustCompile(`<!--[\s\S]*?-->`)
	return re.ReplaceAllString(htmlString, "")
}

// important to allow the notspans to be injected, the spaces after removing the comments cause the regexp to fail
func removeWhitespaceBetweenTags(html string) string {
	// Create a regular expression to match spaces between angle brackets
	re := regexp.MustCompile(`>\s+<`)
	// Replace all matches of spaces between angle brackets with "><"
	return re.ReplaceAllString(html, "><")
}

// Function to hash a struct using SHA-256
func hashStruct(s interface{}) ([]byte, error) {
	// Convert struct to JSON
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	// Hash the JSON data using SHA-256
	hasher := sha256.New()
	hasher.Write(jsonData)
	hash := hasher.Sum(nil)

	return hash, nil
}
