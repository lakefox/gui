package flex

import (
	"gui/cstyle"
	"gui/cstyle/plugins/inline"
	"gui/element"
	"gui/utils"
	"sort"
	"strings"
)

// !ISSUES: Text disapearing (i think its the inline plugin)
// + height adjust on wrap
// + full screen positioning issues

func Init() cstyle.Plugin {
	return cstyle.Plugin{
		Selector: func(n *element.Node) bool {
			styles := map[string]string{
				"display": "flex",
			}
			matches := true
			for name, value := range styles {
				if n.Style[name] != value && !(value == "*") && n.Style[name] != "" {
					matches = false
				}
			}
			return matches
		},
		Level: 3,
		Handler: func(n *element.Node, state *map[string]element.State) {
			s := *state
			self := s[n.Properties.Id]

			verbs := strings.Split(n.Style["flex-direction"], "-")
			flexDirection := verbs[0]
			if flexDirection == "" {
				flexDirection = "row"
			}
			flexReversed := false
			if len(verbs) > 1 {
				flexReversed = true
			}

			var flexWrapped bool
			if n.Style["flex-wrap"] == "wrap" {
				flexWrapped = true
			} else {
				flexWrapped = false
			}

			hAlign := n.Style["align-content"]
			if hAlign == "" {
				hAlign = "normal"
			}
			vAlign := n.Style["align-items"]
			if vAlign == "" {
				vAlign = "normal"
			}
			justify := n.Style["justify-items"]
			if justify == "" {
				justify = "normal"
			}
			// fmt.Println(flexDirection, flexReversed, flexWrapped, hAlign, vAlign, justify)

			if flexDirection == "row" {

				// Reverse elements
				if flexReversed {
					flexReverse(n, state)
				}
				// Get inital sizing
				textTotal := 0
				textCounts := []int{}
				widths := []float32{}
				innerSizes := [][]float32{}
				minWidths := []float32{}
				maxWidths := []float32{}
				for _, v := range n.Children {
					count := countText(v)
					textTotal += count
					textCounts = append(textCounts, count)

					minw := getMinWidth(&v, state)
					minWidths = append(minWidths, minw)

					maxw := getMaxWidth(&v, state)
					maxWidths = append(maxWidths, maxw)

					w, h := getInnerSize(&v, state)
					innerSizes = append(innerSizes, []float32{w, h})
				}
				selfWidth := (self.Width - self.Padding.Left) - self.Padding.Right
				// if the elements are less than the size of the parent, don't change widths. Just set mins
				if !flexWrapped {
					if add2d(innerSizes, 0) < selfWidth {
						for i, v := range n.Children {
							vState := s[v.Properties.Id]

							w := innerSizes[i][0]
							w -= vState.Margin.Left + vState.Margin.Right + (vState.Border.Width * 2)
							widths = append(widths, w)
						}
					} else {
						// Modifiy the widths so they aren't under the mins
						for i, v := range n.Children {
							vState := s[v.Properties.Id]

							w := ((selfWidth / float32(textTotal)) * float32(textCounts[i]))
							w -= vState.Margin.Left + vState.Margin.Right + (vState.Border.Width * 2)

							if w < minWidths[i] {
								selfWidth -= minWidths[i] + vState.Margin.Left + vState.Margin.Right + (vState.Border.Width * 2)
								textTotal -= textCounts[i]
								textCounts[i] = 0
							}

						}
						for i, v := range n.Children {
							vState := s[v.Properties.Id]

							w := ((selfWidth / float32(textTotal)) * float32(textCounts[i]))
							w -= vState.Margin.Left + vState.Margin.Right + (vState.Border.Width * 2)
							// (w!=w) is of NaN
							if w < minWidths[i] || (w != w) {
								w = minWidths[i]
							}
							widths = append(widths, w)
						}
					}
					// Apply the new widths
					fState := s[n.Children[0].Properties.Id]
					for i, v := range n.Children {
						vState := s[v.Properties.Id]

						vState.Width = widths[i]
						xStore := vState.X
						if i > 0 {
							sState := s[n.Children[i-1].Properties.Id]
							vState.X = sState.X + sState.Width + sState.Margin.Right + vState.Margin.Left + sState.Border.Width + vState.Border.Width
							propagateOffsets(&v, xStore, vState.Y, vState.X, fState.Y, state)
						}

						vState.Y = fState.Y

						(*state)[v.Properties.Id] = vState
						deInline(&v, state)
						applyInline(&v, state)
						applyBlock(&v, state)
					}

					// Set the heights based on the tallest one
					if n.Style["height"] == "" {

						innerSizes = [][]float32{}
						for _, v := range n.Children {
							w, h := getInnerSize(&v, state)
							innerSizes = append(innerSizes, []float32{w, h})
						}
						sort.Slice(innerSizes, func(i, j int) bool {
							return innerSizes[i][1] > innerSizes[j][1]
						})
					} else {
						innerSizes[0][1] = self.Height
					}
					for _, v := range n.Children {
						vState := s[v.Properties.Id]
						vState.Height = innerSizes[0][1]
						(*state)[v.Properties.Id] = vState
					}
				} else {
					// Flex Wrapped
					sum := innerSizes[0][0]
					shifted := false
					for i := 0; i < len(n.Children); i++ {
						v := n.Children[i]
						vState := s[v.Properties.Id]

						// if the next plus current will break then
						w := innerSizes[i][0]
						if i > 0 {
							sib := s[n.Children[i-1].Properties.Id]
							if w+sum > selfWidth {
								if maxWidths[i] > selfWidth {
									w = selfWidth - vState.Margin.Left - vState.Margin.Right - (vState.Border.Width * 2)
								}
								sum = 0
								shifted = true
							} else {
								if !shifted {
									propagateOffsets(&v, vState.X, vState.Y, vState.X, sib.Y, state)

									vState.Y = sib.Y
									(*state)[v.Properties.Id] = vState
								} else {
									shifted = false
								}
								sum += w
							}
						}

						widths = append(widths, w)
					}

					// Move the elements into the correct position
					rows := [][]int{}
					start := 0
					maxH := float32(0)
					var prevOffset float32
					for i := 0; i < len(n.Children); i++ {
						v := n.Children[i]
						vState := s[v.Properties.Id]

						vState.Width = widths[i]
						xStore := vState.X
						yStore := vState.Y

						if i > 0 {
							sib := s[n.Children[i-1].Properties.Id]
							if vState.Y+prevOffset == sib.Y {
								vState.Y += prevOffset

								if vState.Height < sib.Height {
									vState.Height = sib.Height
								}
								// Shift right if on a row with sibling
								xStore = sib.X + sib.Width + sib.Margin.Right + sib.Border.Width + vState.Margin.Left + vState.Border.Width
							} else {
								// Shift under sibling
								yStore = sib.Y + sib.Height + sib.Margin.Top + sib.Margin.Bottom + sib.Border.Width*2
								prevOffset = yStore - vState.Y
								rows = append(rows, []int{start, i, int(maxH)})
								start = i
								maxH = 0
							}
							propagateOffsets(&v, vState.X, vState.Y, xStore, yStore, state)
						}
						vState.X = xStore
						vState.Y = yStore

						(*state)[v.Properties.Id] = vState
						deInline(&v, state)
						applyInline(&v, state)
						applyBlock(&v, state)
						_, h := getInnerSize(&v, state)
						h = utils.Max(h, vState.Height)
						maxH = utils.Max(maxH, h)
						vState.Height = h
						(*state)[v.Properties.Id] = vState
					}
					if start < len(n.Children)-1 {
						rows = append(rows, []int{start, len(n.Children) - 1, int(maxH)})
					}
					for _, v := range rows {
						for i := v[0]; i < v[1]; i++ {
							vState := s[n.Children[i].Properties.Id]
							vState.Height = float32(v[2])
							(*state)[n.Children[i].Properties.Id] = vState
						}
					}
				}

				// Shift to the right if reversed
				if flexReversed {
					last := s[n.Children[len(n.Children)-1].Properties.Id]
					offset := (self.X + self.Width - self.Padding.Right) - (last.X + last.Width + last.Margin.Right + last.Border.Width)
					for i, v := range n.Children {
						vState := s[v.Properties.Id]
						propagateOffsets(&n.Children[i], vState.X, vState.Y, vState.X+offset, vState.Y, state)
						vState.X += offset

						(*state)[v.Properties.Id] = vState
					}
				}

			}

			// Column doesn't really need a lot done bc it is basically block styling rn
			if flexDirection == "column" && flexReversed {
				flexReverse(n, state)
			}
			if n.Style["height"] == "" {
				_, h := getInnerSize(n, state)
				self.Height = h
			}
			(*state)[n.Properties.Id] = self
		},
	}
}

func applyBlock(n *element.Node, state *map[string]element.State) {
	accum := float32(0)
	inlineOffset := float32(0)
	s := *state
	lastHeight := float32(0)
	baseY := s[n.Children[0].Properties.Id].Y
	for i := 0; i < len(n.Children); i++ {
		v := &n.Children[i]
		vState := s[v.Properties.Id]

		if v.Style["display"] != "block" {
			vState.Y += inlineOffset
			accum = (vState.Y - baseY)
			lastHeight = vState.Height
		} else if v.Style["position"] != "absolute" {
			vState.Y += accum
			inlineOffset += (vState.Height + (vState.Border.Width * 2) + vState.Margin.Top + vState.Margin.Bottom + vState.Padding.Top + vState.Padding.Bottom) + lastHeight
		}
		(*state)[v.Properties.Id] = vState
	}
}

func deInline(n *element.Node, state *map[string]element.State) {
	s := *state
	// self := s[n.Properties.Id]
	baseX := float32(-1)
	baseY := float32(-1)
	for _, v := range n.Children {
		vState := s[v.Properties.Id]

		if v.Style["display"] == "inline" {
			if baseX < 0 && baseY < 0 {
				baseX = vState.X
				baseY = vState.Y
			} else {
				vState.X = baseX
				vState.Y = baseY
				(*state)[v.Properties.Id] = vState

			}
		} else {
			baseX = float32(-1)
			baseY = float32(-1)
		}

		if len(v.Children) > 0 {
			deInline(&v, state)
		}
	}

}

func applyInline(n *element.Node, state *map[string]element.State) {
	pl := inline.Init()
	for i := 0; i < len(n.Children); i++ {
		v := &n.Children[i]

		if len(v.Children) > 0 {
			applyInline(v, state)
		}

		if pl.Selector(v) {
			pl.Handler(v, state)
		}
	}
}

func propagateOffsets(n *element.Node, prevx, prevy, newx, newy float32, state *map[string]element.State) {
	s := *state
	for _, v := range n.Children {
		vState := s[v.Properties.Id]
		xStore := (vState.X - prevx) + newx
		yStore := (vState.Y - prevy) + newy

		if len(v.Children) > 0 {
			propagateOffsets(&v, vState.X, vState.Y, xStore, yStore, state)
		}
		vState.X = xStore
		vState.Y = yStore
		(*state)[v.Properties.Id] = vState
	}

}

func countText(n element.Node) int {
	count := 0
	groups := []int{}
	for _, v := range n.Children {
		if v.TagName == "notaspan" {
			count += 1
		}
		if v.Style["display"] == "block" {
			groups = append(groups, count)
			count = 0
		}
		if len(v.Children) > 0 {
			count += countText(v)
		}
	}
	groups = append(groups, count)

	sort.Slice(groups, func(i, j int) bool {
		return groups[i] > groups[j]
	})
	return groups[0]
}

func getMinWidth(n *element.Node, state *map[string]element.State) float32 {
	s := *state
	self := s[n.Properties.Id]
	selfWidth := float32(0)

	if len(n.Children) > 0 {
		for _, v := range n.Children {
			selfWidth = utils.Max(selfWidth, getNodeWidth(&v, state))
		}
	} else {
		selfWidth = self.Width
	}

	selfWidth += self.Padding.Left + self.Padding.Right
	return selfWidth
}
func getMaxWidth(n *element.Node, state *map[string]element.State) float32 {
	s := *state
	self := s[n.Properties.Id]
	selfWidth := float32(0)

	if len(n.Children) > 0 {
		for _, v := range n.Children {
			selfWidth += getNodeWidth(&v, state)
		}
	} else {
		selfWidth = self.Width
	}

	selfWidth += self.Padding.Left + self.Padding.Right
	return selfWidth
}

func getNodeWidth(n *element.Node, state *map[string]element.State) float32 {
	s := *state
	self := s[n.Properties.Id]
	w := float32(0)
	w += self.Padding.Left
	w += self.Padding.Right

	w += self.Margin.Left
	w += self.Margin.Right

	w += self.Width

	w += self.Border.Width * 2

	for _, v := range n.Children {
		w = utils.Max(w, getNodeWidth(&v, state))
	}

	return w
}

func getInnerSize(n *element.Node, state *map[string]element.State) (float32, float32) {
	s := *state
	self := s[n.Properties.Id]

	minx := float32(10e10)
	maxw := float32(0)
	miny := float32(10e10)
	maxh := float32(0)
	for _, v := range n.Children {
		vState := s[v.Properties.Id]
		minx = utils.Min(vState.X, minx)
		miny = utils.Min(vState.Y, miny)
		hOffset := (vState.Border.Width * 2) + vState.Margin.Top + vState.Margin.Bottom
		wOffset := (vState.Border.Width * 2) + vState.Margin.Left + vState.Margin.Right
		maxw = utils.Max(vState.X+vState.Width+wOffset, maxw)
		maxh = utils.Max(vState.Y+vState.Height+hOffset, maxh)
	}
	w := maxw - minx
	h := maxh - miny

	// !ISSUE: this is a hack to get things moving adding 13 is random
	w += self.Padding.Left + self.Padding.Right + 13
	h += self.Padding.Top + self.Padding.Bottom
	return w, h
}

func add2d(arr [][]float32, index int) float32 {
	var sum float32
	if len(arr) == 0 {
		return sum
	}

	for i := 0; i < len(arr); i++ {
		if len(arr[i]) <= index {
			return sum
		}
		sum += arr[i][index]
	}

	return sum
}

func flexReverse(n *element.Node, state *map[string]element.State) {
	s := *state
	tempNodes := []element.Node{}
	tempStates := []element.State{}
	for i := len(n.Children) - 1; i >= 0; i-- {
		tempNodes = append(tempNodes, n.Children[i])
		tempStates = append(tempStates, s[n.Children[i].Properties.Id])
	}

	for i := 0; i < len(tempStates); i++ {
		vState := s[n.Children[i].Properties.Id]
		propagateOffsets(&n.Children[i], vState.X, vState.Y, vState.X, tempStates[i].Y, state)
		vState.Y = tempStates[i].Y
		(*state)[n.Children[i].Properties.Id] = vState
	}

	n.Children = tempNodes
}
