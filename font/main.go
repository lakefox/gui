package font

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// LoadSystemFont loads a font from the system fonts directory or loads a specific font by name
func GetFontPath(fontName string, bold, italic bool) string {

	if len(fontName) == 0 {
		fontName = "serif"
	}

	// Check if a special font family is requested
	switch fontName {
	case "sans-serif":
		return tryLoadSystemFont("Arial", bold, italic)
	case "monospace":
		return tryLoadSystemFont("Andle Mono", bold, italic)
	case "serif":
		return tryLoadSystemFont("Georgia", bold, italic)
	}

	// Use the default font if the specified font is not found
	return tryLoadSystemFont(fontName, bold, italic)
}

var allFonts, _ = getSystemFonts()

func tryLoadSystemFont(fontName string, bold, italic bool) string {
	font := fontName
	if bold {
		font += " Bold"
	}
	if italic {
		font += " Italic"
	}
	for _, v := range allFonts {
		if strings.Contains(v, "/"+font) {
			return v
		}
	}

	return ""
}

func sortByLength(strings []string) {
	sort.Slice(strings, func(i, j int) bool {
		return len(strings[i]) < len(strings[j])
	})
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

func LoadFont(fontName string, fontSize int, bold, italic bool) (font.Face, error) {
	// Use a TrueType font file for the specified font name
	fontFile := GetFontPath(fontName, bold, italic)
	println(fontFile)

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

	options := truetype.Options{
		Size:    float64(fontSize),
		DPI:     72,
		Hinting: font.HintingNone,
	}

	// Create a new font face with the specified size
	return truetype.NewFace(fnt, &options), nil
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

func getSystemFonts() ([]string, error) {
	var fontPaths []string

	switch runtime.GOOS {
	case "windows":
		fontPaths = append(fontPaths, getWindowsFontPaths()...)
	case "darwin":
		fontPaths = append(fontPaths, getMacFontPaths()...)
	case "linux":
		fontPaths = append(fontPaths, getLinuxFontPaths()...)
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	sortByLength(fontPaths)

	return fontPaths, nil
}

func getWindowsFontPaths() []string {
	var fontPaths []string

	// System Fonts
	systemFontsDir := "C:\\Windows\\Fonts"
	getFontsRecursively(systemFontsDir, &fontPaths)

	// User Fonts
	userFontsDir := os.ExpandEnv("%APPDATA%\\Microsoft\\Windows\\Fonts")
	getFontsRecursively(userFontsDir, &fontPaths)

	return fontPaths
}

func getMacFontPaths() []string {
	var fontPaths []string

	// System Fonts
	systemFontsDirs := []string{"/System/Library/Fonts", "/Library/Fonts"}
	for _, dir := range systemFontsDirs {
		getFontsRecursively(dir, &fontPaths)
	}

	// User Fonts
	userFontsDir := filepath.Join(os.Getenv("HOME"), "Library/Fonts")
	getFontsRecursively(userFontsDir, &fontPaths)

	return fontPaths
}

func getLinuxFontPaths() []string {
	var fontPaths []string

	// System Fonts
	systemFontsDirs := []string{"/usr/share/fonts", "/usr/local/share/fonts"}
	for _, dir := range systemFontsDirs {
		getFontsRecursively(dir, &fontPaths)
	}

	// User Fonts
	userFontsDir := filepath.Join(os.Getenv("HOME"), ".fonts")
	getFontsRecursively(userFontsDir, &fontPaths)

	return fontPaths
}

func getFontsRecursively(dir string, fontPaths *[]string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			getFontsRecursively(path, fontPaths)
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".ttf") {
			*fontPaths = append(*fontPaths, path)
		}
	}
}

type Text struct {
	Text          string
	Font          font.Face
	Color         color.Color
	Image         *image.RGBA
	Underlined    bool
	LineThrough   bool
	Align         string
	Indent        int
	LetterSpacing int
	LineHeight    int
	WordSpacing   int
	WhiteSpace    string
	Shadows       []Shadow
}

type Shadow struct {
	X     int
	Y     int
	Blur  int
	Color color.Color
}

func (t *Text) Render() {
	width := MeasureText(t.Font, t.Text)
	height := t.Font.Metrics().Height.Round()

	// Use fully transparent color for the background
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	r, g, b, a := t.Color.RGBA()

	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{uint8(r), uint8(g), uint8(b), 0}}, image.Point{}, draw.Over)

	dot := fixed.Point26_6{X: fixed.I(0), Y: t.Font.Metrics().Ascent}

	dr := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}},
		Face: t.Font,
		Dot:  dot,
	}

	dr.DrawString(t.Text)

	t.Image = img
}
