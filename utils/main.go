package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func AddMarginAndPadding(styleMap map[string]map[string]string, id string, width, height float32) (float32, float32) {
	setMP(id, styleMap)
	if styleMap[id]["margin-left"] != "" || styleMap[id]["margin-right"] != "" {
		l, _ := ConvertToPixels(styleMap[id]["margin-left"], width)
		r, _ := ConvertToPixels(styleMap[id]["margin-right"], width)
		width += l
		width += r
	}
	if styleMap[id]["margin-top"] != "" || styleMap[id]["margin-bottom"] != "" {
		t, _ := ConvertToPixels(styleMap[id]["margin-top"], height)
		b, _ := ConvertToPixels(styleMap[id]["margin-bottom"], height)
		height += t
		height += b
	}
	if styleMap[id]["padding-left"] != "" || styleMap[id]["padding-right"] != "" {
		l, _ := ConvertToPixels(styleMap[id]["padding-left"], width)
		r, _ := ConvertToPixels(styleMap[id]["padding-right"], width)
		width += l
		width += r
	}
	if styleMap[id]["padding-top"] != "" || styleMap[id]["padding-bottom"] != "" {
		t, _ := ConvertToPixels(styleMap[id]["padding-top"], height)
		b, _ := ConvertToPixels(styleMap[id]["padding-bottom"], height)
		height += t
		height += b
	}
	return width, height
}

func setMP(id string, styleMap map[string]map[string]string) {
	if styleMap[id]["margin"] != "" {
		left, right, top, bottom := convertMarginToIndividualProperties(styleMap[id]["margin"])
		styleMap[id]["margin-left"] = left
		styleMap[id]["margin-right"] = right
		styleMap[id]["margin-top"] = top
		styleMap[id]["margin-bottom"] = bottom
	}
	if styleMap[id]["padding"] != "" {
		left, right, top, bottom := convertMarginToIndividualProperties(styleMap[id]["margin"])
		styleMap[id]["padding-left"] = left
		styleMap[id]["padding-right"] = right
		styleMap[id]["padding-top"] = top
		styleMap[id]["padding-bottom"] = bottom
	}
}

func convertMarginToIndividualProperties(margin string) (string, string, string, string) {
	// Remove extra whitespace
	margin = strings.TrimSpace(margin)

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
func ConvertToPixels(value string, max float32) (float32, error) {
	// Define conversion factors for different units
	unitFactors := map[string]float32{
		"px": 1,
		"em": 16,    // Assuming 1em = 16px (typical default font size in browsers)
		"pt": 1.33,  // Assuming 1pt = 1.33px (typical conversion)
		"pc": 16.89, // Assuming 1pc = 16.89px (typical conversion)
		"%":  max / 100,
		// "vw": c.Width / 100,
		// "vh": c.Height / 100,
	}

	// Extract numeric value and unit using regular expression
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([a-zA-Z\%]+)$`)
	match := re.FindStringSubmatch(value)

	if len(match) != 3 {
		return 0, fmt.Errorf("invalid input format")
	}

	numericValue, err := (strconv.ParseFloat(match[1], 64))
	numericValue32 := float32(numericValue)
	check(err)

	unit, ok := unitFactors[match[2]]
	if !ok {
		return 0, fmt.Errorf("unsupported unit: %s", match[2])
	}

	return numericValue32 * unit, nil
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
