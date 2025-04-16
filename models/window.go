package models

import rl "github.com/gen2brain/raylib-go/raylib"

// Window struct for UI
type Window struct {
	title      string
	rect       rl.Rectangle
	isDragging bool
	dragOffset rl.Vector2
	content    func(*Window) // Function to render window content
	data       interface{}   // Generic data field for window content
	isActive   bool          // Whether this window is currently active/focused
}
