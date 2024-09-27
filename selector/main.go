package selector

import (
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

	// Check if there's a `::` and split the string
	var prePseudo string
	var pseudo string

	if idx := strings.Index(s, "::"); idx != -1 {
		prePseudo = s[:idx]
		pseudo = s[idx:] // Keep everything after `::` together
	} else {
		prePseudo = s
	}

	// Process the part before the pseudo-element selector
	for _, char := range prePseudo {
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

	// Add the pseudo-element (if any) as a single item
	if pseudo != "" {
		result = append(result, pseudo)
	}

	return result
}

func Contains(selector []string, node []string) bool {
	selectorSet := make(map[string]struct{}, len(node))
	for _, s := range node {
		selectorSet[s] = struct{}{}
	}

	for _, s := range selector {
		if _, exists := selectorSet[s]; !exists {
			return false
		}
	}
	return true
}

// precompile selectors so you don't have to split selector on every style sheet and every selector everytime, put in a map
// if the selector is the same as another merge them.
// if the selector has a parent or child selector then make a css prop "Selector": "full selector". should be blank if not
// can run the actual test in that case, just see if a part of the selector matches, can be cleaned up later
