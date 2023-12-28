package events

import (
	"gui/element"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// need to make the element id branch to allow for quick element finding

type Events struct {
	RightClick  bool
	LeftClick   bool
	Click       bool
	MouseEnter  bool
	MouseExit   bool
	MouseMove   bool
	DoubleClick bool
	MouseOver   bool
	Resize      bool
}

func GetEvents(elements []element.Node) element.Node {
	mp := rl.GetMousePosition()
	// Mouse over
	for i := len(elements) - 1; i >= 0; i-- {
		if elements[i].X < mp.X && elements[i].X+elements[i].Width > mp.X {
			if elements[i].Y < mp.Y && elements[i].Y+elements[i].Height > mp.Y {
				return elements[i]
			}
		}
	}
	return element.Node{}
}
