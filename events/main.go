package events

import (
	"gui/element"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// need to make the element id branch to allow for quick element finding

func GetEvents(el *element.Node) {
	mp := rl.GetMousePosition()
	// Mouse over
	loop(el, mp)
}

func loop(el *element.Node, mp rl.Vector2) {
	if el.X < mp.X && el.X+el.Width > mp.X {
		if el.Y < mp.Y && el.Y+el.Height > mp.Y {
			// Mouse is over element
			// fmt.Println(el.Id, (el.EventListeners))
			evt := element.Event{
				X:     int(mp.X),
				Y:     int(mp.Y),
				Click: false,
			}

			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				evt.Click = true
				if len(el.EventListeners["click"]) > 0 {
					for _, handler := range el.EventListeners["click"] {
						handler(evt)
					}
				}
			}

		}
	}
	for _, v := range el.Children {
		loop(&v, mp)
	}
}
