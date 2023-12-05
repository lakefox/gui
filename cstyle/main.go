package cstyle

import (
	"fmt"
	"gui/parser"
	"os"
	"regexp"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type CSS struct {
	StyleSheets []map[string]map[string]string
}

func (c *CSS) StyleSheet(path string) {
	// Parse the CSS file
	dat, err := os.ReadFile(path)
	check(err)
	styles := parser.ParseCSS(string(dat))

	c.StyleSheets = append(c.StyleSheets, styles)
}

func (c *CSS) StyleTag(css string) {
	styles := parser.ParseCSS(css)
	c.StyleSheets = append(c.StyleSheets, styles)
}

// ConvertToPixels converts a CSS measurement to pixels.
func ConvertToPixels(value string) (float64, error) {
	// Define conversion factors for different units
	unitFactors := map[string]float64{
		"px": 1,
		"em": 16,    // Assuming 1em = 16px (typical default font size in browsers)
		"pt": 1.33,  // Assuming 1pt = 1.33px (typical conversion)
		"pc": 16.89, // Assuming 1pc = 16.89px (typical conversion)
		// Add more units as needed
	}

	// Extract numeric value and unit using regular expression
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([a-zA-Z]+)$`)
	match := re.FindStringSubmatch(value)

	if len(match) != 3 {
		return 0, fmt.Errorf("invalid input format")
	}

	numericValue, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse numeric value: %v", err)
	}

	unit, ok := unitFactors[match[2]]
	if !ok {
		return 0, fmt.Errorf("unsupported unit: %s", match[2])
	}

	return numericValue * unit, nil
}
