package gui

import (
	"bufio"
	_ "embed"
	"gui/cstyle"
	"gui/cstyle/plugins/block"
	"gui/cstyle/plugins/flex"
	"gui/cstyle/plugins/inline"
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
	css.AddPlugin(block.Init())
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

func View(data *Window, width, height int32) {
	data.Document.Properties.Computed["width"] = float32(width)
	data.Document.Properties.Computed["height"] = float32(height)
	data.Document.Style["width"] = strconv.Itoa(int(width)) + "px"
	data.Document.Style["height"] = strconv.Itoa(int(height)) + "px"

	wm := window.NewWindowManager()
	// wm.FPS = true

	wm.OpenWindow(width, height)
	defer wm.CloseWindow()

	evts := map[string]element.EventList{}

	eventStore := &evts

	// Main game loop
	for !wm.WindowShouldClose() {
		rl.BeginDrawing()

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

			data.Document.Style["width"] = strconv.Itoa(int(width)) + "px"
			data.Document.Style["height"] = strconv.Itoa(int(height)) + "px"
			data.Document.Properties.Computed["width"] = float32(width)
			data.Document.Properties.Computed["height"] = float32(height)
		}

		eventStore = events.GetEvents(&data.Document.Children[0], eventStore)
		data.CSS.ComputeNodeStyle(&data.Document.Children[0])
		rd := data.CSS.Render(data.Document.Children[0])
		wm.LoadTextures(rd)
		wm.Draw(rd)

		events.RunEvents(eventStore)

		rl.EndDrawing()
	}
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
		newNode.InnerText = utils.GetInnerText(node)
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

	// println(encapsulateText(htmlContent))
	doc, _ := html.Parse(strings.NewReader(encapsulateText(htmlContent)))

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
	htmlString = removeHTMLComments(htmlString)
	openOpen := regexp.MustCompile(`(<\w+[^>]*>)([^<]+)(<\w+[^>]*>)`)
	closeOpen := regexp.MustCompile(`(</\w+[^>]*>)([^<]+)(<\w+[^>]*>)`)
	closeClose := regexp.MustCompile(`(</\w+[^>]*>)([^<]+)(</\w+[^>]*>)`)
	a := matchFactory(openOpen)
	t := openOpen.ReplaceAllStringFunc(htmlString, a)
	b := matchFactory(closeOpen)
	u := closeOpen.ReplaceAllStringFunc(t, b)
	c := matchFactory(closeClose)
	v := closeClose.ReplaceAllStringFunc(u, c)
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
	re := regexp.MustCompile(`<!--(.*?)-->`)
	return re.ReplaceAllString(htmlString, "")
}
