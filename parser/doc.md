# Parser

```go
package parser

import (
 "regexp"
 "strings"
)

func ParseCSS(css string) map[string]map[string]string {
 selectorMap := make(map[string]map[string]string)

 // Remove comments
 css = removeComments(css)

 // Parse regular selectors and styles
 selectorRegex := regexp.MustCompile(`([^{]+){([^}]+)}`)
 matches := selectorRegex.FindAllStringSubmatch(css, -1)

 for _, match := range matches {
  selectorBlock := strings.TrimSpace(match[1])
  styleBlock := match[2]

  selectors := parseSelectors(selectorBlock)
  for _, selector := range selectors {
   selectorMap[selector] = parseStyles(styleBlock)
  }
 }

 return selectorMap
}

func parseSelectors(selectorBlock string) []string {
 // Split by comma and trim each selector
 selectors := strings.Split(selectorBlock, ",")
 for i, selector := range selectors {
  selectors[i] = strings.TrimSpace(selector)
 }
 return selectors
}

func parseStyles(styleBlock string) map[string]string {
 styleRegex := regexp.MustCompile(`([a-zA-Z-]+)\s*:\s*([^;]+);`)
 matches := styleRegex.FindAllStringSubmatch(styleBlock, -1)

 styleMap := make(map[string]string)
 for _, match := range matches {
  propName := strings.TrimSpace(match[1])
  propValue := strings.TrimSpace(match[2])
  styleMap[propName] = propValue
 }

 return styleMap
}

func ParseStyleAttribute(styleValue string) map[string]string {
 styleMap := make(map[string]string)

 // Regular expression to match key-value pairs in the style attribute
 re := regexp.MustCompile(`\s*([\w-]+)\s*:\s*([^;]+)\s*;`)

 // Find all matches in the style attribute value
 matches := re.FindAllStringSubmatch(styleValue, -1)

 // Populate the map with key-value pairs
 for _, match := range matches {
  if len(match) == 3 {
   key := strings.TrimSpace(match[1])
   value := strings.TrimSpace(match[2])
   styleMap[key] = value
  }
 }

 return styleMap
}

func removeComments(css string) string {
 commentRegex := regexp.MustCompile(`(?s)/\*.*?\*/`)
 return commentRegex.ReplaceAllString(css, "")
}

```
