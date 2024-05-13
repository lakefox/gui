package gui

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"gui/cstyle"
	"gui/cstyle/plugins/flex"
	"gui/cstyle/plugins/inline"
	"gui/cstyle/plugins/inlineText"
	"gui/cstyle/plugins/textAlign"
	"gui/window"

	"gui/element"
	"gui/events"
	"gui/utils"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/net/html"
)

// _ "net/http/pprof"

//go:embed master.css
var mastercss string

type Window struct {
	CSS      cstyle.CSS
	Document element.Node
	Adapter  func()
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
	css.AddPlugin(inlineText.Init())
	css.AddPlugin(flex.Init())

	el := element.Node{}
	document := el.CreateElement("ROOT")
	document.Style["width"] = "800px"
	document.Style["height"] = "450px"
	document.Properties.Id = "ROOT"
	return Window{
		CSS:      css,
		Document: document,
	}
}

func (w *Window) Render(state map[string]element.State) []element.State {
	flatDoc := flatten(w.Document)
	store := []element.State{}

	for _, v := range flatDoc {
		// wont work unless state is a pointer to the og
		// s := state[v.Properties.Id]
		// s.RenderCount++
		// state[v.Properties.Id] = s
		store = append(store, state[v.Properties.Id])
	}
	return store
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

// !ISSUE: Need to make it skip over non modified elements

func View(data *Window, width, height int32) {
	debug := false
	data.Document.Style["width"] = strconv.Itoa(int(width)) + "px"
	data.Document.Style["height"] = strconv.Itoa(int(height)) + "px"

	wm := window.NewWindowManager()
	wm.FPS = true

	wm.OpenWindow(width, height)
	defer wm.CloseWindow()

	evts := map[string]element.EventList{}

	eventStore := &evts

	state := map[string]element.State{}

	state["ROOT"] = element.State{
		Style: map[string]string{
			"width":  strconv.Itoa(int(width)) + "px",
			"height": strconv.Itoa(int(height)) + "px",
		},
	}

	shouldStop := false

	var hash []byte

	// Main game loop
	for !wm.WindowShouldClose() && !shouldStop {
		// fmt.Println("######################")
		rl.BeginDrawing()
		if !shouldStop && debug {
			shouldStop = true
		}
		// Check if the window size has changed
		newWidth := int32(rl.GetScreenWidth())
		newHeight := int32(rl.GetScreenHeight())

		if newWidth != width || newHeight != height {
			rl.ClearBackground(rl.RayWhite)
			// Window has been resized, handle the event
			width = newWidth
			height = newHeight

			data.CSS.Width = float32(width)
			data.CSS.Height = float32(height)

			state["ROOT"].Style["width"] = strconv.Itoa(int(width)) + "px"
			state["ROOT"].Style["height"] = strconv.Itoa(int(height)) + "px"
		}

		newHash, _ := hashStruct(&data.Document.Children[0])
		eventStore = events.GetEvents(&data.Document.Children[0], &state, eventStore)
		if !bytes.Equal(hash, newHash) {
			fmt.Println("rteload")
			hash = newHash
			newDoc := CopyNode(data.Document.Children[0], &data.Document)
			data.CSS.ComputeNodeStyle(&newDoc, &state)
			// !NOTE: Add inner and outerhtml here
		}

		rd := data.Render(state)
		wm.LoadTextures(rd)
		wm.Draw(rd)

		events.RunEvents(eventStore)

		rl.EndDrawing()
	}
}

func CopyNode(node element.Node, parent *element.Node) element.Node {
	n := element.Node{}
	n.TagName = node.TagName
	n.InnerText = node.InnerText
	n.Style = node.Style
	n.Id = node.Id
	n.ClassList = node.ClassList
	n.Href = node.Href
	n.Src = node.Src
	n.Title = node.Title
	n.Attribute = node.Attribute
	n.Value = node.Value
	n.ScrollY = node.ScrollY
	n.InnerHTML = node.InnerHTML
	n.OuterHTML = node.OuterHTML
	n.Properties.Id = node.Properties.Id
	n.Properties.Focusable = node.Properties.Focusable
	n.Properties.Focused = node.Properties.Focused
	n.Properties.Editable = node.Properties.Editable
	n.Properties.Hover = node.Properties.Hover
	n.Properties.Selected = node.Properties.Selected

	n.Parent = parent

	for _, v := range node.Children {
		n.Children = append(n.Children, CopyNode(v, &n))
	}
	return n
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
		parent.AppendChild(newNode)

	} else {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode {
				CreateNode(child, parent)
			}
		}
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
			return submatches[1] + "<text>" + submatches[2] + "</text>" + submatches[3]
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
