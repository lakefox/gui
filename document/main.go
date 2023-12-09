package document

import (
	"bufio"
	"gui/color"
	"gui/cstyle"
	"gui/font"
	"gui/painter"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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

	render(css, d, wm)

	// Main game loop
	for !wm.WindowShouldClose() {
		// Check if the window size has changed
		newWidth := int32(rl.GetScreenWidth())
		newHeight := int32(rl.GetScreenHeight())

		if newWidth != screenWidth || newHeight != screenHeight {
			// Window has been resized, handle the event
			screenWidth = newWidth
			screenHeight = newHeight

			css.Width = float32(screenWidth)
			css.Height = float32(screenHeight)
			render(css, d, wm)
		}

		// Draw rectangles
		wm.DrawRectangles()
	}
}

func render(css cstyle.CSS, d Doc, wm *painter.WindowManager) {
	p := css.Map(d.DOM)

	for _, v := range p.Render {
		styles := p.StyleMap[v.Id]

		x, _ := strconv.ParseFloat(styles["x"], 32)
		y, _ := strconv.ParseFloat(styles["y"], 32)
		width, _ := strconv.ParseFloat(styles["width"], 32)
		height, _ := strconv.ParseFloat(styles["height"], 32)

		if height == 0 {
			continue
		}

		bgColor := color.Background(styles)
		fontColor, _ := color.Font(styles)

		fontFile := font.GetFont(styles)

		var text painter.Text = painter.Text{}

		if len(dom.Children(v.Node)) == 0 {
			text.Value = dom.InnerText(v.Node)
			fs, _ := strconv.ParseFloat(styles["fontSize"], 32)
			if fs == 0 {
				fs = 16
			}
			text.Size = float32(fs)
			text.Color = rl.NewColor(fontColor.R, fontColor.G, fontColor.B, fontColor.A)
			text.Font = fontFile
		}

		node := painter.Rect{
			Node:  rl.NewRectangle(float32(x), float32(y), float32(width), float32(height)),
			Color: rl.NewColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A),
			Text:  text,
		}

		wm.AddRectangle(node)
	}
}

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
