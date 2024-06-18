package selector

import (
	"slices"
	"strings"

	"golang.org/x/net/html"
)

// !ISSUE: Create :not and other selectors

func GetInitCSSSelectors(node *html.Node, selectors []string) []string {
	if node.Type == html.ElementNode {
		selectors = append(selectors, node.Data)
		for _, attr := range node.Attr {
			if attr.Key == "class" {
				classes := strings.Split(attr.Val, " ")
				for _, class := range classes {
					selectors = append(selectors, "."+class)
				}
			} else if attr.Key == "id" {
				selectors = append(selectors, "#"+attr.Val)
			} else {
				selectors = append(selectors, "["+attr.Key+"=\""+attr.Val+"\"]")
			}
		}
	}

	return selectors
}

func SplitSelector(s string) []string {
	var result []string
	var current strings.Builder

	for _, char := range s {
		switch char {
		case '.', '#', '[', ']', ':':
			if current.Len() > 0 {
				if char == ']' {
					current.WriteRune(char)
				}
				result = append(result, current.String())
				current.Reset()
			}
			if char != ']' {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func Contains(selector []string, node []string) bool {
	has := true
	for _, s := range selector {
		if !slices.Contains(node, s) {
			has = false
		}
	}
	return has
}
