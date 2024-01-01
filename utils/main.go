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

func AddMarginAndPadding(styleMap map[string]map[string]string, id string, width, height float32) (float32, float32, float32, float32) {
	fs := font.GetFontSize(styleMap[id])
	if styleMap[id]["padding-left"] != "" || styleMap[id]["padding-right"] != "" {
		l, _ := ConvertToPixels(styleMap[id]["padding-left"], fs, width)
		r, _ := ConvertToPixels(styleMap[id]["padding-right"], fs, width)
		width += l
		width += r
	}
	if styleMap[id]["padding-top"] != "" || styleMap[id]["padding-bottom"] != "" {
		t, _ := ConvertToPixels(styleMap[id]["padding-top"], fs, height)
		b, _ := ConvertToPixels(styleMap[id]["padding-bottom"], fs, height)
		height += t
		height += b
	}

	var marginWidth, marginHeight float32 = width, height

	if styleMap[id]["margin-left"] != "" || styleMap[id]["margin-right"] != "" {
		l, _ := ConvertToPixels(styleMap[id]["margin-left"], fs, width)
		r, _ := ConvertToPixels(styleMap[id]["margin-right"], fs, width)
		marginWidth += l
		marginWidth += r
	}
	if styleMap[id]["margin-top"] != "" || styleMap[id]["margin-bottom"] != "" {
		t, _ := ConvertToPixels(styleMap[id]["margin-top"], fs, height)
		b, _ := ConvertToPixels(styleMap[id]["margin-bottom"], fs, height)
		marginHeight += t
		marginHeight += b
	}
	return width, height, marginWidth, marginHeight
}

func SetMP(id string, styleMap map[string]map[string]string) {
	if styleMap[id]["margin"] != "" {
		left, right, top, bottom := convertMarginToIndividualProperties(styleMap[id]["margin"])
		if styleMap[id]["margin-left"] == "" {
			styleMap[id]["margin-left"] = left
		}
		if styleMap[id]["margin-right"] == "" {
			styleMap[id]["margin-right"] = right
		}
		if styleMap[id]["margin-top"] == "" {
			styleMap[id]["margin-top"] = top
		}
		if styleMap[id]["margin-bottom"] == "" {
			styleMap[id]["margin-bottom"] = bottom
		}
	}
	if styleMap[id]["padding"] != "" {
		left, right, top, bottom := convertMarginToIndividualProperties(styleMap[id]["padding"])
		if styleMap[id]["padding-left"] == "" {
			styleMap[id]["padding-left"] = left
		}
		if styleMap[id]["padding-right"] == "" {
			styleMap[id]["padding-right"] = right
		}
		if styleMap[id]["padding-top"] == "" {
			styleMap[id]["padding-top"] = top
		}
		if styleMap[id]["padding-bottom"] == "" {
			styleMap[id]["padding-bottom"] = bottom
		}
	}
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
	pos := styleMap[n.Id]["position"]

	if pos == "relative" {
		x, _ := strconv.ParseFloat(styleMap[n.Id]["x"], 32)
		y, _ := strconv.ParseFloat(styleMap[n.Id]["y"], 32)
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
		return nil, fmt.Errorf("Expected a pointer to a struct")
	}

	// Get the struct field by name
	field := val.Elem().FieldByName(fieldName)

	// Check if the field exists
	if !field.IsValid() {
		return nil, fmt.Errorf("Field not found: %s", fieldName)
	}

	return field.Interface(), nil
}

func getAttributes(node *html.Node) map[string]string {
	attributes := make(map[string]string)

	for _, attr := range node.Attr {
		attributes[attr.Key] = attr.Val
	}

	return attributes
}

func setAttribute(node *html.Node, key, value string) {
	// Check if the node is an element node
	if node.Type == html.ElementNode {
		// Iterate through the attributes
		for i, attr := range node.Attr {
			// If the attribute key matches, update its value
			if attr.Key == key {
				node.Attr[i].Val = value
				return
			}
		}

		// If the attribute key was not found, add a new attribute
		node.Attr = append(node.Attr, html.Attribute{
			Key: key,
			Val: value,
		})
	}
}
