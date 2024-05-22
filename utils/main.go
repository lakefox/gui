package utils

import (
	"bytes"
	"fmt"
	"gui/element"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func GetXY(n *element.Node, state *map[string]element.State) (float32, float32) {
	s := *state
	// self := s[n.Properties.Id]

	offsetX := float32(0)
	offsetY := float32(0)

	if n.Parent != nil {
		parent := s[n.Parent.Properties.Id]
		// x, y := GetXY(n.Parent, state)
		offsetX += parent.Border.Width + parent.Padding.Left
		offsetY += parent.Border.Width + parent.Padding.Top
	}

	return offsetX, offsetY
}

type WidthHeight struct {
	Width  float32
	Height float32
}

func GetWH(n element.Node, state *map[string]element.State) WidthHeight {
	s := *state
	self := s[n.Properties.Id]
	parent := element.State{}

	if n.Style == nil {
		n.Style = map[string]string{}
	}

	fs := self.EM

	var pwh WidthHeight
	if n.Parent != nil {
		parent = s[n.Parent.Properties.Id]
		pwh = GetWH(*n.Parent, state)
	} else {
		pwh = WidthHeight{}
		if n.Style["width"] != "" {
			str := strings.TrimSuffix(n.Style["width"], "px")
			// Convert the string to float32
			f, _ := strconv.ParseFloat(str, 32)
			pwh.Width = float32(f)
		}

		if n.Style["height"] != "" {
			str := strings.TrimSuffix(n.Style["height"], "px")
			// Convert the string to float32
			f, _ := strconv.ParseFloat(str, 32)
			pwh.Height = float32(f)
		}
	}

	wStyle := n.Style["width"]

	if wStyle == "" && n.Style["display"] != "inline" {
		wStyle = "100%"
	}

	width, _ := ConvertToPixels(wStyle, fs, pwh.Width)
	height, _ := ConvertToPixels(n.Style["height"], fs, pwh.Height)

	if n.Style["min-width"] != "" {
		minWidth, _ := ConvertToPixels(n.Style["min-width"], fs, pwh.Width)
		width = Max(width, minWidth)
	}

	if n.Style["max-width"] != "" {
		maxWidth, _ := ConvertToPixels(n.Style["max-width"], fs, pwh.Width)
		width = Min(width, maxWidth)
	}
	if n.Style["min-height"] != "" {
		minHeight, _ := ConvertToPixels(n.Style["min-height"], fs, pwh.Height)
		height = Max(height, minHeight)
	}

	if n.Style["max-height"] != "" {
		maxHeight, _ := ConvertToPixels(n.Style["max-height"], fs, pwh.Height)
		height = Min(height, maxHeight)
	}

	wh := WidthHeight{
		Width:  width,
		Height: height,
	}

	if n.Parent != nil {
		wh.Width += self.Padding.Left + self.Padding.Right
		wh.Height += self.Padding.Top + self.Padding.Bottom
		// fmt.Println(n.Properties.Id, wh, p)
	}

	if wStyle == "100%" {
		wh.Width = wh.Width - ((self.Margin.Right + self.Margin.Left + (self.Border.Width * 2)) + (parent.Padding.Left + parent.Padding.Right) + (self.Padding.Left + self.Padding.Right))
	}

	if n.Style["height"] == "100%" {
		wh.Height -= (self.Margin.Top + self.Margin.Bottom) + (parent.Padding.Top + parent.Padding.Bottom)
	}

	return wh
}

func GetMP(n element.Node, wh WidthHeight, state *map[string]element.State, t string) element.MarginPadding {
	s := *state
	self := s[n.Properties.Id]
	fs := self.EM
	m := element.MarginPadding{}

	if n.Style[t] != "" {
		left, right, top, bottom := convertMarginToIndividualProperties(n.Style[t])
		if n.Style[t+"-left"] == "" {
			n.Style[t+"-left"] = left
		}
		if n.Style[t+"-right"] == "" {
			n.Style[t+"-right"] = right
		}
		if n.Style[t+"-top"] == "" {
			n.Style[t+"-top"] = top
		}
		if n.Style[t+"-bottom"] == "" {
			n.Style[t+"-bottom"] = bottom
		}
	}
	if n.Style[t+"-left"] != "" || n.Style[t+"-right"] != "" {
		l, _ := ConvertToPixels(n.Style[t+"-left"], fs, wh.Width)
		r, _ := ConvertToPixels(n.Style[t+"-right"], fs, wh.Width)
		m.Left = l
		m.Right = r
	}
	if n.Style[t+"-top"] != "" || n.Style[t+"-bottom"] != "" {
		top, _ := ConvertToPixels(n.Style[t+"-top"], fs, wh.Height)
		b, _ := ConvertToPixels(n.Style[t+"-bottom"], fs, wh.Height)
		m.Top = top
		m.Bottom = b
	}
	if t == "margin" {
		if n.Style["margin"] == "auto" && n.Style["margin-left"] == "" && n.Style["margin-right"] == "" {
			pwh := GetWH(*n.Parent, state)
			m.Left = Max((pwh.Width-wh.Width)/2, 0)
			m.Right = m.Left
		}
	}

	return m
}

func convertMarginToIndividualProperties(margin string) (string, string, string, string) {
	// Remove extra whitespace
	margin = strings.TrimSpace(margin)

	if margin == "" {
		return "0px", "0px", "0px", "0px"
	}

	// Regular expression to match values with optional units
	re := regexp.MustCompile(`(-?\d+(\.\d+)?)(\w*|\%)?`)

	// Extract numerical values from the margin property
	matches := re.FindAllStringSubmatch(margin, -1)

	// Initialize variables for individual margins
	var left, right, top, bottom string

	switch len(matches) {
	case 1:
		// If only one value is provided, apply it to all margins
		left = matches[0][0]
		right = matches[0][0]
		top = matches[0][0]
		bottom = matches[0][0]
	case 2:
		// If two values are provided, apply the first to top and bottom, and the second to left and right
		top = matches[0][0]
		bottom = matches[0][0]
		left = matches[1][0]
		right = matches[1][0]
	case 3:
		// If three values are provided, apply the first to top, the second to left and right, and the third to bottom
		top = matches[0][0]
		left = matches[1][0]
		right = matches[1][0]
		bottom = matches[2][0]
	case 4:
		// If four values are provided, apply them to top, right, bottom, and left, respectively
		top = matches[0][0]
		right = matches[1][0]
		bottom = matches[2][0]
		left = matches[3][0]
	}

	return left, right, top, bottom
}

// ConvertToPixels converts a CSS measurement to pixels.
func ConvertToPixels(value string, em, max float32) (float32, error) {
	unitFactors := map[string]float32{
		"px": 1,
		"em": em,
		"pt": 1.33,
		"pc": 16.89,
		"%":  max / 100,
		"vw": max / 100,
		"vh": max / 100,
		"cm": 37.79527559,
		"in": 96,
	}

	if strings.HasPrefix(value, "calc(") {
		// Handle calculation expression
		calcResult, err := evaluateCalcExpression(value[5:len(value)-1], em, max)
		if err != nil {
			return 0, err
		}
		return calcResult, nil
	} else {
		// Extract numeric value and unit
		for k, v := range unitFactors {
			if strings.HasSuffix(value, k) {
				cutStr, _ := strings.CutSuffix(value, k)
				numericValue, _ := strconv.ParseFloat(cutStr, 64)
				return float32(numericValue) * v, nil
			}
		}
		return 0, fmt.Errorf("unable to parse value")
	}

}

// evaluateCalcExpression recursively evaluates 'calc()' expressions
func evaluateCalcExpression(expression string, em, max float32) (float32, error) {
	terms := strings.FieldsFunc(expression, func(c rune) bool {
		return c == '+' || c == '-' || c == '*' || c == '/'
	})

	operators := strings.FieldsFunc(expression, func(c rune) bool {
		return c != '+' && c != '-' && c != '*' && c != '/'
	})

	var result float32

	for i, term := range terms {
		value, err := ConvertToPixels(strings.TrimSpace(term), em, max)
		if err != nil {
			return 0, err
		}

		if i > 0 {
			switch operators[i-1] {
			case "+":
				result += value
			case "-":
				result -= value
			case "*":
				result *= value
			case "/":
				if value != 0 {
					result /= value
				} else {
					return 0, fmt.Errorf("division by zero in 'calc()' expression")
				}
			}
		} else {
			result = value
		}
	}

	return result, nil
}

func Merge(m1, m2 map[string]string) map[string]string {
	// Create a new map and copy m1 into it
	result := make(map[string]string)
	for k, v := range m1 {
		result[k] = v
	}

	// Merge m2 into the new map
	for k, v := range m2 {
		result[k] = v
	}

	return result
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

// getStructField uses reflection to get the value of a struct field by name
func GetStructField(data interface{}, fieldName string) (interface{}, error) {
	val := reflect.ValueOf(data)

	// Make sure we have a pointer to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a pointer to a struct")
	}

	// Get the struct field by name
	field := val.Elem().FieldByName(fieldName)

	// Check if the field exists
	if !field.IsValid() {
		return nil, fmt.Errorf("field not found: %s", fieldName)
	}

	return field.Interface(), nil
}

func SetStructFieldValue(data interface{}, fieldName string, newValue interface{}) error {
	val := reflect.ValueOf(data)

	// Make sure we have a pointer to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct")
	}

	// Get the struct field by name
	field := val.Elem().FieldByName(fieldName)

	// Check if the field exists
	if !field.IsValid() {
		return fmt.Errorf("field not found: %s", fieldName)
	}

	// Check if the new value type is assignable to the field type
	if !reflect.ValueOf(newValue).Type().AssignableTo(field.Type()) {
		return fmt.Errorf("incompatible types for field %s", fieldName)
	}

	// Set the new value
	field.Set(reflect.ValueOf(newValue))

	return nil
}

// func Check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

func GetInnerText(n *html.Node) string {
	var result strings.Builder

	var getText func(*html.Node)
	getText = func(n *html.Node) {
		// Skip processing if the node is a head tag
		if n.Type == html.ElementNode && n.Data == "head" {
			return
		}

		// If it's a text node, append its content
		if n.Type == html.TextNode {
			result.WriteString(n.Data)
		}

		// Traverse child nodes recursively
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getText(c)
		}
	}

	getText(n)

	return result.String()
}

func GetPositionOffsetNode(n *element.Node) *element.Node {
	pos := n.Style["position"]

	if pos == "relative" {
		return n
	} else {
		if n.Parent.TagName != "ROOT" {
			if n.Parent.Style != nil {
				return GetPositionOffsetNode(n.Parent)
			} else {
				return nil
			}
		} else {
			return n.Parent
		}
	}
}

func IsParent(n element.Node, name string) bool {
	if n.Parent.TagName != "ROOT" {
		if n.Parent.TagName == name {
			return true
		} else {
			return IsParent(*n.Parent, name)
		}
	} else {
		return false
	}
}

func ChildrenHaveText(n *element.Node) bool {
	for _, child := range n.Children {
		if len(strings.TrimSpace(child.InnerText)) != 0 {
			return true
		}
		// Recursively check if any child nodes have text
		if ChildrenHaveText(&child) {
			return true
		}
	}
	return false
}

func NodeToHTML(node element.Node) (string, string) {
	// if node.TagName == "notaspan" {
	// 	return node.InnerText + " ", ""
	// }

	var buffer bytes.Buffer
	buffer.WriteString("<" + node.TagName)

	if node.Properties.Editable {
		buffer.WriteString(" contentEditable=\"true\"")
	}

	// Add ID if present
	if node.Id != "" {
		buffer.WriteString(" id=\"" + node.Id + "\"")
	}

	// Add ID if present
	if node.Title != "" {
		buffer.WriteString(" title=\"" + node.Title + "\"")
	}

	// Add ID if present
	if node.Src != "" {
		buffer.WriteString(" src=\"" + node.Src + "\"")
	}

	// Add ID if present
	if node.Href != "" {
		buffer.WriteString(" href=\"" + node.Href + "\"")
	}

	// Add class list if present
	if len(node.ClassList.Classes) > 0 || node.ClassList.Value != "" {
		classes := ""
		for _, v := range node.ClassList.Classes {
			if string(v[0]) != ":" {
				classes += v + " "
			}
		}
		classes = strings.TrimSpace(classes)
		if len(classes) > 0 {
			buffer.WriteString(" class=\"" + classes + "\"")
		}
	}

	// Add style if present
	if len(node.Style) > 0 {

		style := ""
		for key, value := range node.Style {
			if key != "inlineText" {
				style += key + ":" + value + ";"
			}
		}
		style = strings.TrimSpace(style)

		if len(style) > 0 {
			buffer.WriteString(" style=\"" + style + "\"")
		}
	}

	// Add other attributes if present
	for key, value := range node.Attribute {
		if strings.TrimSpace(value) != "" {
			buffer.WriteString(" " + key + "=\"" + value + "\"")
		}
	}

	buffer.WriteString(">")

	// Add inner text if present
	if node.InnerText != "" && !ChildrenHaveText(&node) {
		buffer.WriteString(node.InnerText)
	}
	return buffer.String(), "</" + node.TagName + ">"
}

func OuterHTML(node element.Node) string {
	var buffer bytes.Buffer

	tag, closing := NodeToHTML(node)

	buffer.WriteString(tag)

	// Recursively add children
	for _, child := range node.Children {
		buffer.WriteString(OuterHTML(child))
	}

	buffer.WriteString(closing)

	return buffer.String()
}

func InnerHTML(node element.Node) string {
	var buffer bytes.Buffer
	// Recursively add children
	for _, child := range node.Children {
		buffer.WriteString(OuterHTML(child))
	}
	return buffer.String()
}
