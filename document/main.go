package document

import (
	"bufio"
	"gui/cstyle"
	"gui/element"
	"gui/events"
	"gui/window"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gui/cstyle/plugins/block"
	"gui/cstyle/plugins/flex"
	"gui/cstyle/plugins/inline"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/net/html"
)

type Window struct {
	StyleSheets []string
	StyleTags   []string
	DOM         *html.Node
	Title       string
}

type Document struct {
	CSS cstyle.CSS
}

func (doc Document) Open(index string, script func(*element.Node)) {
	d := parse(index)

	wm := window.NewWindowManager()
	wm.FPS = true

	// Initialization
	var screenWidth int32 = 800
	var screenHeight int32 = 450

	// Open the window
	wm.OpenWindow(screenWidth, screenHeight)
	defer wm.CloseWindow()

	doc.CSS = cstyle.CSS{
		Width:  800,
		Height: 450,
	}
	doc.CSS.StyleSheet("./master.css")
	// css.AddPlugin(position.Init())
	doc.CSS.AddPlugin(inline.Init())
	doc.CSS.AddPlugin(block.Init())
	doc.CSS.AddPlugin(flex.Init())

	for _, v := range d.StyleSheets {
		doc.CSS.StyleSheet(v)
	}

	for _, v := range d.StyleTags {
		doc.CSS.StyleTag(v)
	}

	nodes := doc.CSS.CreateDocument(d.DOM)
	root := &nodes

	script(root)

	// fmt.Println(nodes.Style)

	evts := map[string]element.EventList{}

	eventStore := &evts

	// Main game loop
	for !wm.WindowShouldClose() {
		rl.BeginDrawing()

		// Check if the window size has changed
		newWidth := int32(rl.GetScreenWidth())
		newHeight := int32(rl.GetScreenHeight())

		if newWidth != screenWidth || newHeight != screenHeight {
			rl.ClearBackground(rl.RayWhite)
			// Window has been resized, handle the event
			screenWidth = newWidth
			screenHeight = newHeight

			doc.CSS.Width = float32(screenWidth)
			doc.CSS.Height = float32(screenHeight)

			nodes = doc.CSS.CreateDocument(d.DOM)
			root = &nodes
			script(root)
		}

		eventStore = events.GetEvents(root, eventStore)
		doc.CSS.ComputeNodeStyle(root)
		rd := doc.CSS.Render(*root)
		wm.LoadTextures(rd)
		wm.Draw(rd)

		events.RunEvents(eventStore)

		rl.EndDrawing()
	}
}

// func (doc Document) CreateElement(t string) element.Node {
// 	n := element.Node{TagName: t, Properties: element.Properties{Node: &html.Node{}}}
// 	cstyle.InitNode(&n, doc.CSS)
// 	return n
// }

func parse(path string) Window {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var htmlContent string

	for scanner.Scan() {
		htmlContent += scanner.Text() + "\n"
	}

	check(scanner.Err())
	// println(encapsulateText(htmlContent))
	doc, err := html.Parse(strings.NewReader(encapsulateText(htmlContent)))
	check(err)

	// Extract stylesheet link tags and style tags
	stylesheets := extractStylesheets(doc, filepath.Dir(path))
	styleTags := extractStyleTags(doc)

	d := Window{
		StyleSheets: stylesheets,
		StyleTags:   styleTags,
		DOM:         doc,
	}

	return d
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
