package document

import (
	"bufio"
	"fmt"
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
}

func Parse(path string) Doc {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var htmlContent string

	for scanner.Scan() {
		htmlContent += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	doc, err := dom.FastParse(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
	}

	// Extract stylesheet link tags and style tags
	stylesheets := extractStylesheets(doc, filepath.Dir(path))
	styleTags := extractStyleTags(doc)

	// Print the results
	fmt.Println("Stylesheet Links:")
	for _, link := range stylesheets {
		fmt.Println(link)
	}

	fmt.Println("\nStyle Tags:")
	for _, style := range styleTags {
		fmt.Println(style)
	}

	d := Doc{
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
