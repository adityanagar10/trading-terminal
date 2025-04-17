package models

import (
	colors "github.com/adityanagar10/trader/constants"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Window struct for UI
type Window struct {
	Title            string
	Rect             rl.Rectangle
	IsDragging       bool
	DragOffset       rl.Vector2
	Content          func(*Window)
	Data             interface{}
	IsActive         bool
	ScrollPosition   float32
	MaxScroll        float32
	IsResizing       bool
	ResizeDir        int
	ResizeHandleSize float32
	Padding          float32
	Font             rl.Font
}

func (w *Window) Update(windows []*Window) {
	mousePos := rl.GetMousePosition()

	// Header/title bar area for dragging
	headerRect := rl.Rectangle{
		X:      w.Rect.X,
		Y:      w.Rect.Y,
		Width:  w.Rect.Width,
		Height: 30,
	}

	// Resize handles
	bottomRightRect := rl.Rectangle{
		X:      w.Rect.X + w.Rect.Width - w.ResizeHandleSize,
		Y:      w.Rect.Y + w.Rect.Height - w.ResizeHandleSize,
		Width:  w.ResizeHandleSize,
		Height: w.ResizeHandleSize,
	}

	rightEdgeRect := rl.Rectangle{
		X:      w.Rect.X + w.Rect.Width - w.ResizeHandleSize/2,
		Y:      w.Rect.Y + 30, // Below header
		Width:  w.ResizeHandleSize,
		Height: w.Rect.Height - 30 - w.ResizeHandleSize,
	}

	bottomEdgeRect := rl.Rectangle{
		X:      w.Rect.X,
		Y:      w.Rect.Y + w.Rect.Height - w.ResizeHandleSize/2,
		Width:  w.Rect.Width - w.ResizeHandleSize,
		Height: w.ResizeHandleSize,
	}

	// Check if clicked on this window to make it active
	if rl.CheckCollisionPointRec(mousePos, w.Rect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		// Make this window active and all others inactive
		for _, win := range windows {
			win.IsActive = false
		}
		w.IsActive = true
	}

	// Handle resizing
	if w.IsResizing {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			switch w.ResizeDir {
			case 1: // Bottom-right corner
				w.Rect.Width = mousePos.X - w.Rect.X
				w.Rect.Height = mousePos.Y - w.Rect.Y
			case 2: // Right edge
				w.Rect.Width = mousePos.X - w.Rect.X
			case 3: // Bottom edge
				w.Rect.Height = mousePos.Y - w.Rect.Y
			}

			// Enforce minimum size
			if w.Rect.Width < 100 {
				w.Rect.Width = 100
			}
			if w.Rect.Height < 100 {
				w.Rect.Height = 100
			}
		} else {
			w.IsResizing = false
			w.ResizeDir = 0
		}
	} else if !w.IsDragging {
		// Check if mouse is over resize handles
		if rl.CheckCollisionPointRec(mousePos, bottomRightRect) {
			rl.SetMouseCursor(rl.MouseCursorResizeNWSE)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				w.IsResizing = true
				w.ResizeDir = 1

				// Make this window active
				for _, win := range windows {
					win.IsActive = false
				}
				w.IsActive = true
			}
		} else if rl.CheckCollisionPointRec(mousePos, rightEdgeRect) {
			rl.SetMouseCursor(rl.MouseCursorResizeEW)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				w.IsResizing = true
				w.ResizeDir = 2

				// Make this window active
				for _, win := range windows {
					win.IsActive = false
				}
				w.IsActive = true
			}
		} else if rl.CheckCollisionPointRec(mousePos, bottomEdgeRect) {
			rl.SetMouseCursor(rl.MouseCursorResizeNS)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				w.IsResizing = true
				w.ResizeDir = 3

				// Make this window active
				for _, win := range windows {
					win.IsActive = false
				}
				w.IsActive = true
			}
		} else {
			rl.SetMouseCursor(rl.MouseCursorDefault)
		}
	}

	// Dragging logic
	if rl.CheckCollisionPointRec(mousePos, headerRect) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			w.IsDragging = true
			w.DragOffset = rl.Vector2{
				X: mousePos.X - w.Rect.X,
				Y: mousePos.Y - w.Rect.Y,
			}
			// Make this window active
			for _, win := range windows {
				win.IsActive = false
			}
			w.IsActive = true
		}
	}

	if w.IsDragging {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			w.Rect.X = mousePos.X - w.DragOffset.X
			w.Rect.Y = mousePos.Y - w.DragOffset.Y
		} else {
			w.IsDragging = false
		}
	}

	// Handle scrolling
	if rl.CheckCollisionPointRec(mousePos, w.Rect) {
		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			w.ScrollPosition -= wheel * 20

			// Clamp scroll position
			if w.ScrollPosition < 0 {
				w.ScrollPosition = 0
			}
			if w.ScrollPosition > w.MaxScroll {
				w.ScrollPosition = w.MaxScroll
			}
		}
	}
}

func (w *Window) Draw() {
	// Draw window background with minimal styling
	rl.DrawRectangleRec(w.Rect, colors.ColorPanelBg)

	// Draw window border with different color if active
	borderColor := colors.ColorBorder
	if w.IsActive {
		borderColor = colors.ColorHighlight
	}

	// Draw minimal border (just a bottom and right edge for a modern look)
	// Top border line
	rl.DrawLine(
		int32(w.Rect.X),
		int32(w.Rect.Y),
		int32(w.Rect.X+w.Rect.Width),
		int32(w.Rect.Y),
		borderColor)
	// Right border line
	rl.DrawLine(
		int32(w.Rect.X+w.Rect.Width),
		int32(w.Rect.Y),
		int32(w.Rect.X+w.Rect.Width),
		int32(w.Rect.Y+w.Rect.Height),
		borderColor)
	// Bottom border line
	rl.DrawLine(
		int32(w.Rect.X),
		int32(w.Rect.Y+w.Rect.Height),
		int32(w.Rect.X+w.Rect.Width),
		int32(w.Rect.Y+w.Rect.Height),
		borderColor)
	// Left border line
	rl.DrawLine(
		int32(w.Rect.X),
		int32(w.Rect.Y),
		int32(w.Rect.X),
		int32(w.Rect.Y+w.Rect.Height),
		borderColor)

	// Draw minimal header - just a line with text
	headerHeight := float32(25)

	// Draw title text with monospaced font at proper scale
	rl.DrawTextEx(
		w.Font,
		w.Title,
		rl.Vector2{X: w.Rect.X + w.Padding, Y: w.Rect.Y + 5},
		18,
		1,
		colors.ColorText,
	)

	// Draw header bottom border line
	rl.DrawLine(
		int32(w.Rect.X),
		int32(w.Rect.Y+headerHeight),
		int32(w.Rect.X+w.Rect.Width),
		int32(w.Rect.Y+headerHeight),
		borderColor)

	// Draw content with scissor mode to keep it within bounds
	if w.Content != nil {
		rl.BeginScissorMode(
			int32(w.Rect.X),
			int32(w.Rect.Y+headerHeight),
			int32(w.Rect.Width),
			int32(w.Rect.Height-headerHeight),
		)
		w.Content(w)
		rl.EndScissorMode()
	}

	// Draw resize handle in bottom-right corner (subtle corner)
	if w.IsActive {
		rl.DrawRectangle(
			int32(w.Rect.X+w.Rect.Width-8),
			int32(w.Rect.Y+w.Rect.Height-8),
			8,
			8,
			rl.NewColor(75, 75, 75, 100), // Very subtle handle
		)
	}
}
