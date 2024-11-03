package main

import (
	"fmt"

	"github.com/GeekRicardo/mousehook"
)

func main() {
	mousehook.SetMouseDownCallback(func(event mousehook.MouseEvent) {
		fmt.Printf("Mouse button down: %s at (%d, %d)\n", event.Button, event.X, event.Y)
	})

	mousehook.SetMouseUpCallback(func(event mousehook.MouseEvent) {
		fmt.Printf("Mouse button up: %s at (%d, %d)\n", event.Button, event.X, event.Y)
	})

	mousehook.SetMouseWheelCallback(func(event mousehook.MouseEvent) {
		fmt.Printf("Mouse wheel event: %s -> Delta %d at (%d, %d)\n", event.Button, event.Delta, event.X, event.Y)
	})

	fmt.Println("Press Ctrl+C to exit")
	mousehook.Start()
}
