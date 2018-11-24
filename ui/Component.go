package ui

import (
	"image"
	"image/draw"
)

// Component is the interface type for all UI Components
type Component interface {
	Draw(draw.Image)
	Delete()
}

// BaseComponent is a stub Component made to be inherited from
type BaseComponent struct {
	Position image.Point
	Size     image.Point
}

// Delete frees all resources owned by the Component
func (c *BaseComponent) Delete() {
}

// GetPosition returns the Component's current position
func (c *BaseComponent) GetPosition() image.Point {
	return c.Position
}

// SetPosition sets the Component's position
func (c *BaseComponent) SetPosition(pos image.Point) {
	c.Position = pos
}

// GetSize returns the Component's current size
func (c *BaseComponent) GetSize() image.Point {
	return c.Size
}

// SetSize sets the Component's size
func (c *BaseComponent) SetSize(size image.Point) {
	c.Size = size
}

// Draw renders a BaseComponent
func (c *BaseComponent) Draw(draw.Image) {
}
