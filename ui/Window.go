package ui

import (
	"image"
	"image/draw"
)

type Window struct {
	BaseComponent

	ImgNW     *Image
	ImgNE     *Image
	ImgSW     *Image
	ImgSE     *Image
	ImgN      *Image
	ImgS      *Image
	ImgW      *Image
	ImgE      *Image
	ImgCenter *Image
}

func NewWindow(nw, ne, sw, se, n, s, w, e, c *Image) *Window {
	return &Window{
		ImgNW:     nw,
		ImgNE:     ne,
		ImgSW:     sw,
		ImgSE:     se,
		ImgN:      n,
		ImgS:      s,
		ImgW:      w,
		ImgE:      e,
		ImgCenter: c,
	}
}

func (c *Window) SetPosition(pos image.Point) {
	c.BaseComponent.SetPosition(pos)
	c.updateImages()
}

func (c *Window) SetSize(size image.Point) {
	c.BaseComponent.SetSize(size)
	c.updateImages()
}

func (c *Window) updateImages() {
	pos := c.GetPosition()
	size := c.GetSize()

	width := func(i *Image) int {
		return i.GetSize().X
	}

	height := func(i *Image) int {
		return i.GetSize().Y
	}

	c.ImgNW.SetPosition(pos)
	c.ImgNE.SetPosition(image.Pt(pos.X+size.X-width(c.ImgNE), pos.Y))
	c.ImgSW.SetPosition(image.Pt(pos.X, pos.Y+size.Y-height(c.ImgSW)))
	c.ImgSE.SetPosition(image.Pt(pos.X+size.X-width(c.ImgSE), pos.Y+size.Y-height(c.ImgSE))) //

	c.ImgN.SetPosition(image.Pt(pos.X+width(c.ImgNW), pos.Y))
	c.ImgN.SetSize(image.Pt(size.X-width(c.ImgNW)-width(c.ImgNE), height(c.ImgN)))

	c.ImgS.SetPosition(image.Pt(pos.X+width(c.ImgSW), pos.Y+size.Y-height(c.ImgS)))
	c.ImgS.SetSize(image.Pt(size.X-width(c.ImgSW)-width(c.ImgSE), height(c.ImgS)))

	c.ImgW.SetPosition(image.Pt(pos.X, pos.Y+height(c.ImgNW)))
	c.ImgW.SetSize(image.Pt(width(c.ImgW), size.Y-height(c.ImgNW)-height(c.ImgSW)))

	c.ImgE.SetPosition(image.Pt(pos.X+size.X-width(c.ImgE), pos.Y+height(c.ImgNW)))
	c.ImgE.SetSize(image.Pt(width(c.ImgE), size.Y-height(c.ImgNE)-height(c.ImgSE)))

	c.ImgCenter.SetPosition(pos.Add(c.ImgNW.GetSize()))
	c.ImgCenter.SetSize(size.Sub(c.ImgNW.GetSize()).Sub(c.ImgSE.GetSize()))
}

func (c *Window) Draw(buffer draw.Image) {
	c.ImgNW.Draw(buffer)
	c.ImgNE.Draw(buffer)
	c.ImgSW.Draw(buffer)
	c.ImgSE.Draw(buffer)
	c.ImgN.Draw(buffer)
	c.ImgS.Draw(buffer)
	c.ImgW.Draw(buffer)
	c.ImgE.Draw(buffer)
	c.ImgCenter.Draw(buffer)
}
