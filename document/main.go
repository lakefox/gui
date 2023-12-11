package document

import (
	"bufio"
	"gui/cstyle"
	"gui/painter"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

type Doc struct {
	StyleSheets []string
	StyleTags   []string
	DOM         *html.Node
	Title       string
}

func Open(index string) {
	d := Parse(index)

	wm := painter.NewWindowManager()
	wm.FPS = true

	// Initialization
	var screenWidth int32 = 800
	var screenHeight int32 = 450

	// Open the window
	wm.OpenWindow(d.Title, screenWidth, screenHeight)
	defer wm.CloseWindow()

	css := cstyle.CSS{
		Width:  800,
		Height: 450,
	}
	css.StyleSheet("./master.css")

	for _, v := range d.StyleSheets {
		css.StyleSheet(v)
	}

	for _, v := range d.StyleTags {
		css.StyleTag(v)
	}

	css.Map(d.DOM)

	// Main game loop
	for !wm.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		// Check if the window size has changed
		newWidth := int32(rl.GetScreenWidth())
		newHeight := int32(rl.GetScreenHeight())

		if newWidth != screenWidth || newHeight != screenHeight {
			// Window has been resized, handle the event
			screenWidth = newWidth
			screenHeight = newHeight

			css.Width = float32(screenWidth)
			css.Height = float32(screenHeight)
			// styles := css.Map(d.DOM)
			css.Map(d.DOM)
			// parent := cstyle.Node{
			// 	X:      0,
			// 	Y:      0,
			// 	Width:  float32(screenWidth),
			// 	Height: float32(screenHeight),
			// }
			// render(styles, styles.Document, wm, parent)
		}

		// Draw rectangles
		wm.Draw()

		rl.EndDrawing()
	}
}

// func render(m cstyle.Mapped, n *html.Node, wm *painter.WindowManager, parent cstyle.Node) {
// 	id := dom.GetAttribute(n, "DOMNODEID")
// 	fs := font.GetFontSize(m.StyleMap[id])

// 	x, y := m.Position(n, parent)

// 	width, _ := utils.ConvertToPixels(m.InheritProp(n, "width"), fs, float32(parent.Width))
// 	height, _ := utils.ConvertToPixels(m.InheritProp(n, "height"), fs, float32(parent.Height))

// }

func Parse(path string) Doc {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var htmlContent string

	for scanner.Scan() {
		htmlContent += scanner.Text() + "\n"
	}

	check(scanner.Err())

	doc, err := dom.FastParse(strings.NewReader(htmlContent))
	check(err)

	// Extract stylesheet link tags and style tags
	stylesheets := extractStylesheets(doc, filepath.Dir(path))
	styleTags := extractStyleTags(doc)

	d := Doc{
		StyleSheets: stylesheets,
		StyleTags:   styleTags,
		DOM:         doc,
		Title:       dom.InnerText(dom.GetElementsByTagName(doc, "title")[0]),
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
