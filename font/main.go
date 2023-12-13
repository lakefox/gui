package font

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"runtime"
	"strconv"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// LoadSystemFont loads a font from the system fonts directory or loads a specific font by name
func GetFontPath(fontName string) string {

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
	return tryLoadSystemFont(fontName) + ".ttf"
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
		if css["font-size"][fL-2:] == "em" {
			fs64, _ := strconv.ParseFloat(css["font-size"][0:fL-2], 32)
			fs = float32(fs64)
		}
	}

	return fs
}

func LoadFont(fontName string, fontSize int) (font.Face, error) {
	// Use a TrueType font file for the specified font name
	fontFile := GetFontPath(fontName)

	// Read the font file
	fontData, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}

	// Parse the TrueType font data
	fnt, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	// Create a new font face with the specified size
	return truetype.NewFace(fnt, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     72,
		Hinting: font.HintingNone,
	}), nil
}

func MeasureText(face font.Face, text string) int {
	dot := fixed.Point26_6{}
	var width fixed.Int26_6

	for _, runeValue := range text {
		adv, ok := face.GlyphAdvance(runeValue)
		if !ok {
			continue
		}

		// Calculate the glyph bounds
		bounds, _, _ := face.GlyphBounds(runeValue)

		// Update the total width with the glyph advance and bounds
		width += adv + bounds.Max.X - bounds.Min.X
		dot.X += adv
	}

	return width.Round()
}

type Text struct {
	Text  string
	Font  font.Face
	Color color.Color
	Image *image.RGBA
}

func (t *Text) Render() {
	width := MeasureText(t.Font, t.Text)
	height := t.Font.Metrics().Height.Round()

	// Use fully transparent color for the background
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Over)

	dot := fixed.Point26_6{X: fixed.I(0), Y: t.Font.Metrics().Ascent}
	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: t.Font,
		Dot:  dot,
	}

	dr.DrawString(t.Text)

	t.Image = img
}
