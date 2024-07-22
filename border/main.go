package border

import (
	"gui/canvas"
	"gui/color"
	"gui/element"
	"gui/utils"
	ic "image/color"
	"strings"
)

func Parse(cssProperties map[string]string, self, parent element.State) (element.Border, error) {
	// Define default values
	defaultWidth := "0px"
	defaultStyle := "solid"
	defaultColor := "#000000"
	defaultRadius := "0px"

	// Helper function to parse border component
	parseBorderComponent := func(value string) (width, style, color string) {
		components := strings.Fields(value)
		width, style, color = defaultWidth, defaultStyle, defaultColor
		widthSuffixes := []string{"px", "em", "pt", "pc", "%", "vw", "vh", "cm", "in"}

		for _, component := range components {
			if isWidthComponent(component, widthSuffixes) {
				width = component
			} else {
				switch component {
				case "thin", "medium", "thick":
					width = component
				case "none", "hidden", "dotted", "dashed", "solid", "double", "groove", "ridge", "inset", "outset":
					style = component
				default:
					color = component
				}
			}
		}

		return
	}

	// Helper function to parse border radius component
	parseBorderRadiusComponent := func(value string) float32 {
		if value == "" {
			value = defaultRadius
		}
		return utils.ConvertToPixels(value, self.EM, parent.Width)
	}

	// Parse individual border sides
	topWidth, topStyle, topColor := parseBorderComponent(cssProperties["border-top"])
	rightWidth, rightStyle, rightColor := parseBorderComponent(cssProperties["border-right"])
	bottomWidth, bottomStyle, bottomColor := parseBorderComponent(cssProperties["border-bottom"])
	leftWidth, leftStyle, leftColor := parseBorderComponent(cssProperties["border-left"])

	// Parse shorthand border property
	if border, exists := cssProperties["border"]; exists {
		width, style, color := parseBorderComponent(border)
		if _, exists := cssProperties["border-top"]; !exists {
			topWidth, topStyle, topColor = width, style, color
		}
		if _, exists := cssProperties["border-right"]; !exists {
			rightWidth, rightStyle, rightColor = width, style, color
		}
		if _, exists := cssProperties["border-bottom"]; !exists {
			bottomWidth, bottomStyle, bottomColor = width, style, color
		}
		if _, exists := cssProperties["border-left"]; !exists {
			leftWidth, leftStyle, leftColor = width, style, color
		}
	}

	var topLeftRadius,
		topRightRadius,
		bottomLeftRadius,
		bottomRightRadius float32

	if cssProperties["border-radius"] != "" {
		rad := parseBorderRadiusComponent(cssProperties["border-radius"])
		topLeftRadius = rad
		topRightRadius = rad
		bottomLeftRadius = rad
		bottomRightRadius = rad
	}

	// Parse border-radius
	if cssProperties["border-top-left-radius"] != "" {
		topLeftRadius = parseBorderRadiusComponent(cssProperties["border-top-left-radius"])
	}
	if cssProperties["border-top-right-radius"] != "" {
		topRightRadius = parseBorderRadiusComponent(cssProperties["border-top-right-radius"])
	}
	if cssProperties["border-bottom-left-radius"] != "" {
		bottomLeftRadius = parseBorderRadiusComponent(cssProperties["border-bottom-left-radius"])
	}
	if cssProperties["border-bottom-right-radius"] != "" {
		bottomRightRadius = parseBorderRadiusComponent(cssProperties["border-bottom-right-radius"])
	}

	// Convert to pixels
	topWidthPx := utils.ConvertToPixels(topWidth, self.EM, parent.Width)
	rightWidthPx := utils.ConvertToPixels(rightWidth, self.EM, parent.Width)
	bottomWidthPx := utils.ConvertToPixels(bottomWidth, self.EM, parent.Width)
	leftWidthPx := utils.ConvertToPixels(leftWidth, self.EM, parent.Width)

	// Parse colors
	topParsedColor, _ := color.Color(topColor)
	rightParsedColor, _ := color.Color(rightColor)
	bottomParsedColor, _ := color.Color(bottomColor)
	leftParsedColor, _ := color.Color(leftColor)

	return element.Border{
		Top: element.BorderSide{
			Width: topWidthPx,
			Style: topStyle,
			Color: topParsedColor,
		},
		Right: element.BorderSide{
			Width: rightWidthPx,
			Style: rightStyle,
			Color: rightParsedColor,
		},
		Bottom: element.BorderSide{
			Width: bottomWidthPx,
			Style: bottomStyle,
			Color: bottomParsedColor,
		},
		Left: element.BorderSide{
			Width: leftWidthPx,
			Style: leftStyle,
			Color: leftParsedColor,
		},
		Radius: element.BorderRadius{
			TopLeft:     topLeftRadius,
			TopRight:    topRightRadius,
			BottomLeft:  bottomLeftRadius,
			BottomRight: bottomRightRadius,
		},
	}, nil
}

func Draw(n *element.State) {
	ctx := canvas.NewCanvas(int(n.X+
		n.Width+n.Border.Left.Width+n.Border.Right.Width),
		int(n.Y+n.Height+n.Border.Top.Width+n.Border.Bottom.Width))
	ctx.StrokeStyle = ic.RGBA{0, 0, 0, 255}
	if n.Border.Top.Width > 0 {
		drawBorderSide(ctx, "top", n.Border.Top, n)
	}
	if n.Border.Right.Width > 0 {
		drawBorderSide(ctx, "right", n.Border.Right, n)
	}
	if n.Border.Bottom.Width > 0 {
		drawBorderSide(ctx, "bottom", n.Border.Bottom, n)
	}
	if n.Border.Left.Width > 0 {
		drawBorderSide(ctx, "left", n.Border.Left, n)
	}
	n.Canvas = ctx

}

func drawBorderSide(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	switch border.Style {
	case "solid":
		drawSolidBorder(ctx, side, border, s)
	case "dashed":
		drawDashedBorder(ctx, side, border, s)
	case "dotted":
		drawDottedBorder(ctx, side, border, s)
	case "double":
		drawDoubleBorder(ctx, side, border, s)
	case "groove":
		drawGrooveBorder(ctx, side, border, s)
	case "ridge":
		drawRidgeBorder(ctx, side, border, s)
	case "inset":
		drawInsetBorder(ctx, side, border, s)
	case "outset":
		drawOutsetBorder(ctx, side, border, s)
	default:
		drawSolidBorder(ctx, side, border, s)
	}
}

// Helper function to determine if a component is a width value
func isWidthComponent(component string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(component, suffix) {
			return true
		}
	}
	return false
}

func drawSolidBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	radius := int(s.Border.Radius.TopLeft) // Using one radius for simplicity, adjust as needed
	ctx.StrokeStyle = border.Color
	ctx.LineWidth = float64(border.Width)

	switch side {
	case "top":
		ctx.RoundedRect(int(s.X), int(s.Y), int(s.Width), int(border.Width), radius)
	case "right":
		ctx.RoundedRect(int(s.X+s.Width-border.Width), int(s.Y), int(border.Width), int(s.Height), radius)
	case "bottom":
		ctx.RoundedRect(int(s.X), int(s.Y+s.Height-border.Width), int(s.Width), int(border.Width), radius)
	case "left":
		ctx.RoundedRect(int(s.X), int(s.Y), int(border.Width), int(s.Height), radius)
	}
	ctx.Stroke()
}

func drawDashedBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// dashLength := 10
	// gapLength := 5
	// drawDashLine := func(x1, y1, x2, y2 int) {
	// 	for i := 0; i < int(border.Width); i++ {
	// 		for j := x1; j < x2; j += dashLength + gapLength {
	// 			drawLineHelper(img, j, y1+i, j+dashLength, y2+i, border.Color)
	// 		}
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawDashLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y)
	// case "right":
	// 	drawDashLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y)
	// case "bottom":
	// 	drawDashLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width))
	// case "left":
	// 	drawDashLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y)
	// }
}

func drawDottedBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// dotSize := int(border.Width)
	// drawDot := func(x, y int) {
	// 	for i := 0; i < dotSize; i++ {
	// 		for j := 0; j < dotSize; j++ {
	// 			img.Set(x+i, y+j, border.Color)
	// 		}
	// 	}
	// }
	// switch side {
	// case "top":
	// 	for i := rect.Min.X; i < rect.Max.X; i += dotSize * 2 {
	// 		drawDot(i, rect.Min.Y)
	// 	}
	// case "right":
	// 	for i := rect.Min.Y; i < rect.Max.Y; i += dotSize * 2 {
	// 		drawDot(rect.Max.X-dotSize, i)
	// 	}
	// case "bottom":
	// 	for i := rect.Min.X; i < rect.Max.X; i += dotSize * 2 {
	// 		drawDot(i, rect.Max.Y-dotSize)
	// 	}
	// case "left":
	// 	for i := rect.Min.Y; i < rect.Max.Y; i += dotSize * 2 {
	// 		drawDot(rect.Min.X, i)
	// 	}
	// }
}

func drawDoubleBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// innerOffset := int(border.Width / 3)
	// outerOffset := innerOffset * 2
	// drawLine := func(x1, y1, x2, y2, offset int) {
	// 	for i := 0; i < innerOffset; i++ {
	// 		drawLineHelper(img, x1, y1+i+offset, x2, y2+i+offset, border.Color)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, 0)
	// 	drawLine(rect.Min.X, rect.Min.Y+outerOffset, rect.Max.X, rect.Min.Y+outerOffset, 0)
	// case "right":
	// 	drawLine(rect.Max.X-innerOffset, rect.Min.Y, rect.Max.X-innerOffset, rect.Max.Y, 0)
	// 	drawLine(rect.Max.X-outerOffset, rect.Min.Y, rect.Max.X-outerOffset, rect.Max.Y, 0)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-innerOffset, rect.Max.X, rect.Max.Y-innerOffset, 0)
	// 	drawLine(rect.Min.X, rect.Max.Y-outerOffset, rect.Max.X, rect.Max.Y-outerOffset, 0)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, 0)
	// 	drawLine(rect.Min.X+outerOffset, rect.Min.Y, rect.Min.X+outerOffset, rect.Max.Y, 0)
	// }
}

func drawGrooveBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// shadowColor := ic.RGBA{border.Color.R / 2, border.Color.G / 2, border.Color.B / 2, border.Color.A}
	// highlightColor := ic.RGBA{border.Color.R * 2 / 3, border.Color.G * 2 / 3, border.Color.B * 2 / 3, border.Color.A}
	// drawLine := func(x1, y1, x2, y2 int, col ic.RGBA) {
	// 	for i := 0; i < int(border.Width/2); i++ {
	// 		drawLineHelper(img, x1, y1+i, x2, y2+i, col)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, shadowColor)
	// 	drawLine(rect.Min.X, rect.Min.Y+int(border.Width/2), rect.Max.X, rect.Min.Y+int(border.Width/2), highlightColor)
	// case "right":
	// 	drawLine(rect.Max.X-int(border.Width/2), rect.Min.Y, rect.Max.X-int(border.Width/2), rect.Max.Y, shadowColor)
	// 	drawLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y, highlightColor)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width/2), rect.Max.X, rect.Max.Y-int(border.Width/2), highlightColor)
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width), shadowColor)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, highlightColor)
	// 	drawLine(rect.Min.X+int(border.Width/2), rect.Min.Y, rect.Min.X+int(border.Width/2), rect.Max.Y, shadowColor)
	// }
}

func drawRidgeBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// shadowColor := ic.RGBA{border.Color.R / 2, border.Color.G / 2, border.Color.B / 2, border.Color.A}
	// highlightColor := ic.RGBA{border.Color.R * 2 / 3, border.Color.G * 2 / 3, border.Color.B * 2 / 3, border.Color.A}
	// drawLine := func(x1, y1, x2, y2 int, col ic.RGBA) {
	// 	for i := 0; i < int(border.Width/2); i++ {
	// 		drawLineHelper(img, x1, y1+i, x2, y2+i, col)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, highlightColor)
	// 	drawLine(rect.Min.X, rect.Min.Y+int(border.Width/2), rect.Max.X, rect.Min.Y+int(border.Width/2), shadowColor)
	// case "right":
	// 	drawLine(rect.Max.X-int(border.Width/2), rect.Min.Y, rect.Max.X-int(border.Width/2), rect.Max.Y, highlightColor)
	// 	drawLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y, shadowColor)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width/2), rect.Max.X, rect.Max.Y-int(border.Width/2), shadowColor)
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width), highlightColor)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, shadowColor)
	// 	drawLine(rect.Min.X+int(border.Width/2), rect.Min.Y, rect.Min.X+int(border.Width/2), rect.Max.Y, highlightColor)
	// }
}

func drawInsetBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// shadowColor := ic.RGBA{border.Color.R / 2, border.Color.G / 2, border.Color.B / 2, border.Color.A}
	// highlightColor := ic.RGBA{border.Color.R * 2 / 3, border.Color.G * 2 / 3, border.Color.B * 2 / 3, border.Color.A}
	// drawLine := func(x1, y1, x2, y2 int, col ic.RGBA) {
	// 	for i := 0; i < int(border.Width/2); i++ {
	// 		drawLineHelper(img, x1, y1+i, x2, y2+i, col)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, shadowColor)
	// 	drawLine(rect.Min.X, rect.Min.Y+int(border.Width/2), rect.Max.X, rect.Min.Y+int(border.Width/2), highlightColor)
	// case "right":
	// 	drawLine(rect.Max.X-int(border.Width/2), rect.Min.Y, rect.Max.X-int(border.Width/2), rect.Max.Y, shadowColor)
	// 	drawLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y, highlightColor)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width/2), rect.Max.X, rect.Max.Y-int(border.Width/2), highlightColor)
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width), shadowColor)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, highlightColor)
	// 	drawLine(rect.Min.X+int(border.Width/2), rect.Min.Y, rect.Min.X+int(border.Width/2), rect.Max.Y, shadowColor)
	// }
}

func drawOutsetBorder(ctx *canvas.Canvas, side string, border element.BorderSide, s *element.State) {
	// shadowColor := ic.RGBA{border.Color.R * 2 / 3, border.Color.G * 2 / 3, border.Color.B * 2 / 3, border.Color.A}
	// highlightColor := ic.RGBA{border.Color.R / 2, border.Color.G / 2, border.Color.B / 2, border.Color.A}
	// drawLine := func(x1, y1, x2, y2 int, col ic.RGBA) {
	// 	for i := 0; i < int(border.Width/2); i++ {
	// 		drawLineHelper(img, x1, y1+i, x2, y2+i, col)
	// 	}
	// }
	// switch side {
	// case "top":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, highlightColor)
	// 	drawLine(rect.Min.X, rect.Min.Y+int(border.Width/2), rect.Max.X, rect.Min.Y+int(border.Width/2), shadowColor)
	// case "right":
	// 	drawLine(rect.Max.X-int(border.Width/2), rect.Min.Y, rect.Max.X-int(border.Width/2), rect.Max.Y, highlightColor)
	// 	drawLine(rect.Max.X-int(border.Width), rect.Min.Y, rect.Max.X-int(border.Width), rect.Max.Y, shadowColor)
	// case "bottom":
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width/2), rect.Max.X, rect.Max.Y-int(border.Width/2), shadowColor)
	// 	drawLine(rect.Min.X, rect.Max.Y-int(border.Width), rect.Max.X, rect.Max.Y-int(border.Width), highlightColor)
	// case "left":
	// 	drawLine(rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, shadowColor)
	// 	drawLine(rect.Min.X+int(border.Width/2), rect.Min.Y, rect.Min.X+int(border.Width/2), rect.Max.Y, highlightColor)
	// }
}
