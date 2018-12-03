package ui

import (
	"github.com/WhoBrokeTheBuild/TelcomSim/context"
	"github.com/go-gl/mathgl/mgl32"
)

// Component is the interface type for all UI Components
type Component interface {
	Draw(*context.Render)
	Delete()
}

// BaseComponent is a stub Component made to be inherited from
type BaseComponent struct {
	Position mgl32.Vec2
	Size     mgl32.Vec2
}

// Delete frees all resources owned by the Component
func (c *BaseComponent) Delete() {
}

// GetPosition returns the Component's current position
func (c *BaseComponent) GetPosition() mgl32.Vec2 {
	return c.Position
}

// SetPosition sets the Component's position
func (c *BaseComponent) SetPosition(pos mgl32.Vec2) {
	c.Position = pos
}

// GetSize returns the Component's current size
func (c *BaseComponent) GetSize() mgl32.Vec2 {
	return c.Size
}

// SetSize sets the Component's size
func (c *BaseComponent) SetSize(size mgl32.Vec2) {
	c.Size = size
}

// Draw renders a BaseComponent
func (c *BaseComponent) Draw(ctx *context.Render) {
}
