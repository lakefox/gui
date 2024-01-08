package selector

import (
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func GetCSSSelectors(node *html.Node, selectors []string) []string {
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
	var current string

	for _, char := range s {
		switch char {
		case '.', '#', '[', ']':
			if current != "" {
				if string(char) == "]" {
					current += string(char)
				}
				result = append(result, current)
			}
			current = ""
			if string(char) != "]" {
				current += string(char)
			}
		default:
			current += string(char)
		}
	}

	if current != "" && current != "]" {
		result = append(result, current)
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

// func Query(selector string, n *Node) bool {
// 	parts := strings.Split(selector, ">")

// 	selectors := getCSSSelectors(n.Node, []string{})

// 	part := splitSelector(strings.TrimSpace(parts[len(parts)-1]))

// 	fmt.Println(part, selectors)

// 	has := contains(part, selectors)

// 	if len(parts) == 1 || !has {
// 		return has
// 	} else {
// 		return Query(strings.Join(parts[0:len(parts)-1], ">"), n.Parent)
// 	}
// }

// func main() {
// 	selector := "div.class#id[attr=\"value\"] > div"

// 	node := &html.Node{
// 		Type: html.ElementNode,
// 		Data: "div",
// 		Attr: []html.Attribute{
// 			{Key: "class", Val: "class"},
// 			{Key: "id", Val: "id"},
// 			{Key: "attr", Val: "value"},
// 		},
// 	}

// 	nodeparent := &html.Node{
// 		Type: html.ElementNode,
// 		Data: "div",
// 		Attr: []html.Attribute{
// 			{Key: "class", Val: "class"},
// 			{Key: "id", Val: "id"},
// 			{Key: "attr", Val: "value"},
// 		},
// 	}

// 	n := Node{
// 		Node: node,
// 		Parent: &Node{
// 			Node: nodeparent,
// 		},
// 	}

// 	fmt.Println(Query(selector, &n))
// }
