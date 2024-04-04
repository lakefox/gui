package utils

import (
	"fmt"
	"gui/element"
	"gui/font"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func ComputeStyleString(n element.Node) string {
	styles := map[string]string{}
	wh := GetWH(n)
	styles["width"] = fmt.Sprint(wh.Width)
	styles["width"] = fmt.Sprint(wh.Height)
	m := GetMP(n, "marign")
	styles["margin-top"] = fmt.Sprint(m.Top)
	styles["margin-left"] = fmt.Sprint(m.Left)
	styles["margin-right"] = fmt.Sprint(m.Right)
	styles["margin-bottom"] = fmt.Sprint(m.Bottom)
	p := GetMP(n, "padding")
	styles["padding-top"] = fmt.Sprint(p.Top)
	styles["padding-left"] = fmt.Sprint(p.Left)
	styles["padding-right"] = fmt.Sprint(p.Right)
	styles["padding-bottom"] = fmt.Sprint(p.Bottom)
	return MapToInline(styles)
}

func ComputeStyleMap(n element.Node) map[string]float32 {
	inline := InlineToMap(n.Style["computed"])
	parsed := map[string]float32{}
	mt, _ := strconv.ParseFloat(inline["margin-top"], 32)
	parsed["margin-top"] = float32(mt)
	ml, _ := strconv.ParseFloat(inline["margin-left"], 32)
	parsed["margin-left"] = float32(ml)
	mr, _ := strconv.ParseFloat(inline["margin-right"], 32)
	parsed["margin-right"] = float32(mr)
	mb, _ := strconv.ParseFloat(inline["margin-bottom"], 32)
	parsed["margin-bottom"] = float32(mb)

	pt, _ := strconv.ParseFloat(inline["padding-top"], 32)
	parsed["padding-top"] = float32(pt)
	pl, _ := strconv.ParseFloat(inline["padding-left"], 32)
	parsed["padding-left"] = float32(pl)
	pr, _ := strconv.ParseFloat(inline["padding-right"], 32)
	parsed["padding-right"] = float32(pr)
	pb, _ := strconv.ParseFloat(inline["padding-bottom"], 32)
	parsed["padding-bottom"] = float32(pb)

	width, _ := strconv.ParseFloat(inline["width"], 32)
	parsed["width"] = float32(width)
	height, _ := strconv.ParseFloat(inline["height"], 32)
	parsed["height"] = float32(height)
	return parsed
}

// MapToInlineCSS converts a map[string]string to a string formatted like inline CSS style
func MapToInline(m map[string]string) string {
	var cssStrings []string
	for key, value := range m {
		cssStrings = append(cssStrings, fmt.Sprintf("%s: %s;", key, value))
	}
	return strings.Join(cssStrings, " ")
}

// InlineCSSToMap converts a string formatted like inline CSS style to a map[string]string
func InlineToMap(cssString string) map[string]string {
	cssMap := make(map[string]string)
	declarations := strings.Split(cssString, ";")
	for _, declaration := range declarations {
		parts := strings.Split(strings.TrimSpace(declaration), ":")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			cssMap[key] = value
		}
	}
	return cssMap
}

type WidthHeight struct {
	Width  float32
	Height float32
}

func GetWH(n element.Node) WidthHeight {
	fs := font.GetFontSize(n.Style)

	var pwh WidthHeight
	if n.Parent != nil {
		pwh = GetWH(*n.Parent)
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

	width, _ := ConvertToPixels(n.Style["width"], fs, pwh.Width)
	if n.Style["min-width"] != "" {
		minWidth, _ := ConvertToPixels(n.Style["min-width"], fs, pwh.Width)
		width = Max(width, minWidth)
	}

	if n.Style["max-width"] != "" {
		maxWidth, _ := ConvertToPixels(n.Style["max-width"], fs, pwh.Width)
		width = Min(width, maxWidth)
	}

	height, _ := ConvertToPixels(n.Style["height"], fs, pwh.Height)
	if n.Style["min-height"] != "" {
		minHeight, _ := ConvertToPixels(n.Style["min-height"], fs, pwh.Height)
		height = Max(height, minHeight)
	}

	if n.Style["max-height"] != "" {
		maxHeight, _ := ConvertToPixels(n.Style["max-height"], fs, pwh.Height)
		height = Min(height, maxHeight)
	}
	return WidthHeight{
		Width:  width,
		Height: height,
	}
}

func SetWH(width, height float32) {
	// could have a calculated style so map[string]string{"computed":"width: 100;x: 100;etc..",}
	// then make a function that can parse/update all of them as needed..
	// might still run into the issue of things not always being updated
}

type MarginPadding struct {
	Top    float32
	Left   float32
	Right  float32
	Bottom float32
}

func GetMP(n element.Node, t string) MarginPadding {
	fs := font.GetFontSize(n.Style)
	m := MarginPadding{}

	wh := GetWH(n)

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
			// this dont work
			m.Left = Max((n.Parent.Properties.Width-wh.Width)/2, 0)
			m.Right = Max((n.Parent.Properties.Width-wh.Width)/2, 0)
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
	}

	re := regexp.MustCompile(`calc\(([^)]*)\)|^(\d+(?:\.\d+)?)\s*([a-zA-Z\%]+)$`)
	match := re.FindStringSubmatch(value)

	if match != nil {
		if len(match[1]) > 0 {
			calcResult, err := evaluateCalcExpression(match[1], em, max)
			if err != nil {
				return 0, err
			}
			return calcResult, nil
		}

		if len(match[2]) > 0 && len(match[3]) > 0 {
			numericValue, err := strconv.ParseFloat(match[2], 64)
			if err != nil {
				return 0, fmt.Errorf("error parsing numeric value: %v", err)
			}
			return float32(numericValue) * unitFactors[match[3]], nil
		}
	}

	return 0, fmt.Errorf("invalid input format: %s", value)
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

func GetTextBounds(text string, fontSize, width, height float32) (float32, float32) {
	w := float32(len(text) * int(fontSize))
	h := fontSize
	if width > 0 && height > 0 {
		if w > width {
			height = Max(height, float32(math.Ceil(float64(w/width)))*h)
		}
		return width, height
	} else {
		return w, h
	}

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

func ExMerge(m1, m2 map[string]string) map[string]string {
	// Create a new map and copy m1 into it
	result := make(map[string]string)
	for k, v := range m1 {
		result[k] = v
	}

	// Merge m2 into the new map only if the key is not already present
	for k, v := range m2 {
		if result[k] == "" {
			result[k] = v
		}
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

func FindRelative(n *element.Node, styleMap map[string]map[string]string) (float32, float32) {
	pos := styleMap[n.Properties.Id]["position"]

	if pos == "relative" {
		x, _ := strconv.ParseFloat(styleMap[n.Properties.Id]["x"], 32)
		y, _ := strconv.ParseFloat(styleMap[n.Properties.Id]["y"], 32)
		return float32(x), float32(y)
	} else {
		if n.Parent != nil {
			x, y := FindRelative(n.Parent, styleMap)
			return x, y
		} else {
			return 0, 0
		}
	}
}

func ParseFloat(str string, def float32) float32 {
	var a float32
	if str == "" {
		a = 0
	} else {
		v, _ := strconv.ParseFloat(str, 32)
		a = float32(v)
	}
	return a
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

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetInnerText(n *html.Node) string {
	var result strings.Builder

	var getText func(*html.Node)
	getText = func(n *html.Node) {
		if n.Type == html.TextNode {
			result.WriteString(n.Data)
		}

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
