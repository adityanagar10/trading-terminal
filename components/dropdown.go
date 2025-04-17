package components

import (
	colors "github.com/adityanagar10/trader/constants"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Dropdown represents a dropdown menu component
type Dropdown struct {
	Rect            rl.Rectangle
	Options         []string
	SelectedIndex   int
	IsOpen          bool
	Label           string
	Font            rl.Font
	OnChangeHandler func(index int)
}

// NewDropdown creates a new dropdown with default values
func NewDropdown(x, y, width float32, options []string, label string, font rl.Font) *Dropdown {
	return &Dropdown{
		Rect:          rl.NewRectangle(x, y, width, 26),
		Options:       options,
		SelectedIndex: 0,
		IsOpen:        false,
		Label:         label,
		Font:          font,
	}
}

// SetOnChangeHandler sets the callback for when selection changes
func (d *Dropdown) SetOnChangeHandler(handler func(int)) {
	d.OnChangeHandler = handler
}

// GetSelectedOption returns the currently selected option as a string

func (d *Dropdown) GetSelectedOption() string {
	return d.Options[d.SelectedIndex]
}

func (d *Dropdown) Update() {
	mousePos := rl.GetMousePosition()

	// Check if clicked on dropdown
	if rl.CheckCollisionPointRec(mousePos, d.Rect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		d.IsOpen = !d.IsOpen
	} else if d.IsOpen {
		// Check if clicked on an option
		for i := range d.Options {
			optionRect := rl.Rectangle{
				X:      d.Rect.X,
				Y:      d.Rect.Y + d.Rect.Height + float32(i)*26,
				Width:  d.Rect.Width,
				Height: 26,
			}

			if rl.CheckCollisionPointRec(mousePos, optionRect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				if d.SelectedIndex != i {
					d.SelectedIndex = i
					if d.OnChangeHandler != nil {
						d.OnChangeHandler(i)
					}
				}
				d.IsOpen = false
				break
			}
		}

		// Close if clicked elsewhere
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			d.IsOpen = false
		}
	}
}

func (d *Dropdown) Draw() {
	// Draw main dropdown box
	rl.DrawRectangleRec(d.Rect, colors.ColorPanelBg)
	rl.DrawRectangleLinesEx(d.Rect, 1, colors.ColorBorder)

	// Draw selected option
	selectedOption := d.Options[d.SelectedIndex]
	rl.DrawTextEx(
		d.Font,
		selectedOption,
		rl.Vector2{
			X: d.Rect.X + 10,
			Y: d.Rect.Y + 5,
		},
		16,
		1,
		colors.ColorText)

	// Draw dropdown arrow
	arrowChar := "▼"
	if d.IsOpen {
		arrowChar = "▲"
	}
	rl.DrawText(
		arrowChar,
		int32(d.Rect.X+d.Rect.Width-20),
		int32(d.Rect.Y+5),
		16,
		colors.ColorSubtext)

	// Draw label above dropdown
	rl.DrawTextEx(
		d.Font,
		d.Label,
		rl.Vector2{
			X: d.Rect.X,
			Y: d.Rect.Y - 20,
		},
		16,
		1,
		colors.ColorSubtext)

	// Draw options if open
	if d.IsOpen {
		for i, option := range d.Options {
			optionRect := rl.Rectangle{
				X:      d.Rect.X,
				Y:      d.Rect.Y + d.Rect.Height + float32(i)*26,
				Width:  d.Rect.Width,
				Height: 26,
			}

			// Highlight selected option
			bgColor := colors.ColorPanelBg
			if i == d.SelectedIndex {
				bgColor = colors.ColorHeaderBg
			}

			rl.DrawRectangleRec(optionRect, bgColor)
			rl.DrawRectangleLinesEx(optionRect, 1, colors.ColorBorder)

			rl.DrawTextEx(
				d.Font,
				option,
				rl.Vector2{
					X: optionRect.X + 10,
					Y: optionRect.Y + 5,
				},
				16,
				1,
				colors.ColorText)
		}
	}
}
