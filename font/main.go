package font

import (
	"fmt"
	"gui/element"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

// func MeasureLine(n *element.Node, state *element.State) (int, int) {
// 	passed := false
// 	lineOffset, nodeOffset := 0, 0
// 	for _, v := range n.Parent.Children {
// 		l := MeasureText(state, v.InnerText)
// 		if v.Properties.Id == n.Properties.Id {
// 			passed = true
// 			lineOffset += l
// 		} else {
// 			if !passed {
// 				nodeOffset += l
// 			}
// 			lineOffset += l
// 		}
// 	}
// 	return lineOffset, nodeOffset
// }

func MeasureText(s *element.State, text string) int {
	t := s.Text
	var width fixed.Int26_6

	for _, runeValue := range text {
		if runeValue == ' ' {
			// Handle spaces separately, add word spacing
			width += fixed.I(t.WordSpacing)
		} else {
			adv, ok := t.Font.GlyphAdvance(runeValue)
			if !ok {
				continue
			}

			// Update the total width with the glyph advance and bounds
			width += adv + fixed.I(t.LetterSpacing)
		}
	}

	return width.Round()
}

func MeasureSpace(t *element.Text) int {
	adv, _ := t.Font.GlyphAdvance(' ')
	return adv.Round()
}

func MeasureLongest(s *element.State) int {
	lines := getLines(s)
	var longestLine string
	maxLength := 0

	for _, line := range lines {
		length := len(line)
		if length > maxLength {
			maxLength = length
			longestLine = line
		}
	}
	return MeasureText(s, longestLine)
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

func Render(s *element.State) float32 {
	t := &s.Text
	lines := getLines(s)

	if t.LineHeight == 0 {
		t.LineHeight = t.EM + 3
	}
	// Use fully transparent color for the background
	img := image.NewRGBA(image.Rect(0, 0, t.Width, t.LineHeight*(len(lines))))

	// fmt.Println(t.Width, t.LineHeight, (len(lines)))

	r, g, b, a := t.Color.RGBA()

	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{uint8(r), uint8(g), uint8(b), uint8(0)}}, image.Point{}, draw.Over)
	// fmt.Println(int(t.Font.Metrics().Ascent))
	dot := fixed.Point26_6{X: fixed.I(0), Y: (fixed.I(t.LineHeight+(t.EM/2)) / 2)}

	dr := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}},
		Face: t.Font,
		Dot:  dot,
	}
	t.Image = img

	fh := fixed.I(t.LineHeight)

	for _, v := range lines {
		lineWidth := MeasureText(s, v)
		if t.Align == "justify" {
			dr.Dot.X = 0
			spaces := strings.Count(v, " ")
			if spaces > 1 {
				spacing := fixed.I((t.Width - MeasureText(s, v)) / spaces)

				if spacing > 0 {
					for _, word := range strings.Fields(v) {
						dr.DrawString(word)
						dr.Dot.X += spacing
					}
				} else {
					dr.Dot.X = 0
					drawString(*t, dr, v, lineWidth)
				}
			} else {
				dr.Dot.X = 0
				drawString(*t, dr, v, lineWidth)
			}

		} else {
			if t.Align == "left" || t.Align == "" {
				dr.Dot.X = 0
			} else if t.Align == "center" {
				dr.Dot.X = fixed.I((t.Width - MeasureText(s, v)) / 2)
			} else if t.Align == "right" {
				dr.Dot.X = fixed.I(t.Width - MeasureText(s, v))
			}
			// dr.Dot.X = 0
			drawString(*t, dr, v, lineWidth)
		}
		dr.Dot.Y += fh
	}
	s.Text.X = MeasureText(s, lines[len(lines)-1])
	return float32(t.LineHeight * len(lines))
}

func drawString(t element.Text, dr *font.Drawer, v string, lineWidth int) {
	underlinePosition := dr.Dot
	for _, ch := range v {
		if ch == ' ' {
			// Handle spaces separately, add word spacing
			dr.Dot.X += fixed.I(t.WordSpacing)
		} else {
			dr.DrawString(string(ch))
			dr.Dot.X += fixed.I(t.LetterSpacing)
		}
	}
	if t.Underlined || t.Overlined || t.LineThrough {

		underlinePosition.X = 0
		baseLineY := underlinePosition.Y

		if t.Underlined {
			underlinePosition.Y = baseLineY + t.Font.Metrics().Descent
			drawLine(t.Image, underlinePosition, fixed.Int26_6(lineWidth), t.DecorationThickness, t.DecorationColor)
		}
		if t.LineThrough {
			underlinePosition.Y = baseLineY - (t.Font.Metrics().Descent)
			drawLine(t.Image, underlinePosition, fixed.Int26_6(lineWidth), t.DecorationThickness, t.DecorationColor)
		}
		if t.Overlined {
			underlinePosition.Y = baseLineY - t.Font.Metrics().Descent*3
			drawLine(t.Image, underlinePosition, fixed.Int26_6(lineWidth), t.DecorationThickness, t.DecorationColor)
		}
	}
}

func wrap(s *element.State, breaker string, breakNewLines bool) []string {
	var start int = 0
	strngs := []string{}
	var text []string
	broken := strings.Split(s.Text.Text, breaker)
	re := regexp.MustCompile(`[\r\n]+`)
	if breakNewLines {
		for _, v := range broken {
			text = append(text, re.Split(v, -1)...)
		}
	} else {
		text = append(text, broken...)
	}
	for i := 0; i < len(text); i++ {
		text[i] = re.ReplaceAllString(text[i], "")
	}
	for i := 0; i < len(text); i++ {
		seg := strings.Join(text[start:int(Min(float32(i+1), float32(len(text))))], breaker)
		if MeasureText(s, seg) > s.Text.Width {
			strngs = append(strngs, strings.Join(text[start:i], breaker))
			start = i
		}
	}
	if len(strngs) > 0 {
		strngs = append(strngs, strings.Join(text[start:], breaker))
	} else {
		strngs = append(strngs, strings.Join(text[start:], breaker))
	}
	return strngs
}

func drawLine(img draw.Image, start fixed.Point26_6, width fixed.Int26_6, thickness int, col color.Color) {
	// Bresenham's line algorithm
	x0, y0 := start.X.Round(), start.Y.Round()
	x1 := x0 + int(width)
	y1 := y0
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1

	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}

	err := dx - dy

	for {
		for i := 0; i < thickness; i++ {
			img.Set(x0, (y0-(thickness/2))+i, col)
		}

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getLines(s *element.State) []string {
	t := s.Text
	text := s.Text.Text
	var lines []string
	if t.WhiteSpace == "nowrap" {
		re := regexp.MustCompile(`\s+`)
		s.Text.Text = re.ReplaceAllString(text, " ")
		lines = wrap(s, "<br />", false)
	} else {
		if t.WhiteSpace == "pre" {
			re := regexp.MustCompile("\t")
			s.Text.Text = re.ReplaceAllString(text, "     ")
			nl := regexp.MustCompile(`[\r\n]+`)
			lines = nl.Split(text, -1)
		} else if t.WhiteSpace == "pre-line" {
			re := regexp.MustCompile(`\s+`)
			s.Text.Text = re.ReplaceAllString(text, " ")
			lines = wrap(s, " ", true)
		} else if t.WhiteSpace == "pre-wrap" {
			lines = wrap(s, " ", true)
		} else {
			re := regexp.MustCompile(`\s+`)
			s.Text.Text = re.ReplaceAllString(text, " ")
			nl := regexp.MustCompile(`[\r\n]+`)
			s.Text.Text = nl.ReplaceAllString(text, "")
			// n.InnerText = strings.TrimSpace(text)
			lines = wrap(s, t.WordBreak, false)
		}
		for i, v := range lines {
			lines[i] = v + t.WordBreak
		}
	}
	return lines
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}
