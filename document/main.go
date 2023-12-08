package document

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

type Doc struct {
	StyleSheets []string
	StyleTags   []string
	DOM         *html.Node
	Title       string
}

func Write(path string) Doc {
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
