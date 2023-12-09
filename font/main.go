package font

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

// LoadSystemFont loads a font from the system fonts directory or loads a specific font by name
func GetFont(styleMap map[string]string) string {

	fontName := styleMap["font-family"]

	if len(fontName) == 0 {
		fontName = "serif"
	}

	// Check if a special font family is requested
	switch fontName {
	case "sans-serif":
		return tryLoadSystemFont("Arial.ttf")
	case "monospace":
		return tryLoadSystemFont("Andle Mono.ttf")
	case "serif":
		return tryLoadSystemFont("Georgia.ttf")
	}

	// Use the default font if the specified font is not found
	return tryLoadSystemFont(fontName + ".ttf")
}

func tryLoadSystemFont(fontName string) string {

	// Get the system font directory
	fontDir, err := getSystemFontDir()
	if err != nil {
		fmt.Println("Error getting system font directory:", err)
		return ""
	}

	// Check if the font file exists in the system font directory
	fontPath := fontDir + string(os.PathSeparator) + fontName
	if _, err := os.Stat(fontPath); err != nil {
		fmt.Println("Error checking font file:", err)
		return ""
	}

	return fontPath
}

func getSystemFontDir() (string, error) {
	var fontDir string

	switch runtime.GOOS {
	case "windows":
		fontDir = os.Getenv("SystemRoot") + `\Fonts`
	case "darwin":
		fontDir = "/System/Library/Fonts/Supplemental"
	default:
		fontDir = "/usr/share/fonts/truetype"
	}

	return fontDir, nil
}

func GetFontSize(css map[string]string) float32 {
	fL := len(css["font-size"])

	var fs float32 = 16

	if fL > 0 {
		if css["font-size"][fL-2:] == "px" {
			fs64, _ := strconv.ParseFloat(css["font-size"][0:fL-2], 32)
			fs = float32(fs64)
		}
	}

	return fs
}
