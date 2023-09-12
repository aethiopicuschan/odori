package ui

import (
	"image/color"
	"math"

	"github.com/aethiopicuschan/odori/constant"
	"github.com/aethiopicuschan/odori/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type scrollBar struct {
	show       bool
	height     int
	width      int
	pos        int
	cursorOn   bool
	grabbed    bool
	grabbedAtY int
}

type Explorer struct {
	onPick           func(s sprite.Sprite)
	sprites          []sprite.Sprite
	cursolOn         int
	clicked          int
	doubleClickCount int
	scrollBar        scrollBar
	scrollOffset     float64
	height           int
	totalHeight      int
	offsetX          int
	offsetY          int
	size             int
}

func NewExplorer(onPick func(s sprite.Sprite)) *Explorer {
	_, h := ebiten.WindowSize()
	return &Explorer{
		onPick: onPick,
		sprites: []sprite.Sprite{
			sprite.NewEmptySprite(),
		},
		cursolOn:         -1,
		clicked:          -1,
		doubleClickCount: -1,
		height:           h / 3,
		totalHeight:      0,
		offsetX:          10,
		offsetY:          10,
		size:             100,
		scrollBar: scrollBar{
			show:     false,
			height:   0,
			width:    10,
			pos:      0,
			cursorOn: false,
		},
	}
}

func (e *Explorer) AppendSprite(sprite sprite.Sprite) {
	e.sprites = append(e.sprites, sprite)
}

func (e *Explorer) Update() error {
	if e.doubleClickCount >= 0 {
		e.doubleClickCount++
		// ダブルクリックの判定は1/3秒
		maxDoubleClickCount := ebiten.TPS() / 3
		if e.doubleClickCount > maxDoubleClickCount {
			e.doubleClickCount = -1
			e.clicked = -1
		}
	}
	w, _ := ebiten.WindowSize()
	width := w - constant.MenuWidth
	spritesPerRow := (width - e.offsetX) / (e.size + e.offsetX)
	numOfColumn := int(math.Ceil(float64(len(e.sprites)) / float64(spritesPerRow)))
	e.totalHeight = numOfColumn*(e.size+e.offsetY) + e.offsetY

	if !ebiten.IsFocused() {
		e.scrollBar.cursorOn = false
		e.scrollBar.grabbed = false
		return nil
	}
	cursorX, cursorY := ebiten.CursorPosition()
	isCursorOnExplorer := cursorX >= w-width && cursorX <= w && cursorY >= 0 && cursorY <= e.height
	shape := ebiten.CursorShapeDefault
	if isCursorOnExplorer {
		defer func() {
			ebiten.SetCursorShape(shape)
		}()
	}

	// which sprite is cursor on?
	e.cursolOn = -1
	if len(e.sprites) > 0 && isCursorOnExplorer {
		zeroX := constant.MenuWidth
		xOnExplorer := cursorX - zeroX
		yOnExplorer := cursorY - int(e.scrollOffset)
		row := xOnExplorer / (e.size + e.offsetX)
		if row >= spritesPerRow {
			row = spritesPerRow - 1
		}
		col := yOnExplorer / (e.size + e.offsetY)
		index := row + col*spritesPerRow
		if index < len(e.sprites) {
			targetRectX := zeroX + e.offsetX*(row+1) + (e.size * row)
			targetRectY := e.offsetY*(col+1) + (e.size * col) + int(e.scrollOffset)
			if cursorX >= targetRectX && cursorX <= targetRectX+e.size && cursorY >= targetRectY && cursorY <= targetRectY+e.size {
				e.cursolOn = index
				shape = ebiten.CursorShapePointer
				// handle click
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					if e.doubleClickCount == -1 {
						e.doubleClickCount = 0
						e.clicked = e.cursolOn
					} else {
						e.doubleClickCount = -1
						if e.cursolOn == e.clicked {
							e.onPick(e.sprites[e.cursolOn])
						}
					}
				}
			}
		}
	}

	// Logic of ScrollBar.
	if e.totalHeight <= e.height {
		e.scrollBar.show = false
		return nil
	}
	e.scrollBar.show = true
	e.scrollBar.height = int(float64(e.height) * (float64(e.height) / float64(e.totalHeight)))
	if isCursorOnExplorer {
		_, wheelY := ebiten.Wheel()
		e.scrollBar.pos -= int(wheelY * 2.0)
		if e.scrollBar.pos < 0 {
			e.scrollBar.pos = 0
		}
		if e.scrollBar.pos > e.height-e.scrollBar.height {
			e.scrollBar.pos = e.height - e.scrollBar.height
		}
	}

	if cursorX >= w-e.scrollBar.width && cursorX <= w && cursorY >= e.scrollBar.pos && cursorY <= e.scrollBar.pos+e.scrollBar.height {
		e.scrollBar.cursorOn = true
		shape = ebiten.CursorShapePointer
	} else {
		if e.scrollBar.cursorOn {
			e.scrollBar.cursorOn = false
		}
	}
	if e.scrollBar.cursorOn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		e.scrollBar.grabbed = true
		e.scrollBar.grabbedAtY = cursorY
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		e.scrollBar.grabbed = false
	}
	if e.scrollBar.grabbed {
		e.scrollBar.pos += cursorY - e.scrollBar.grabbedAtY
		e.scrollBar.grabbedAtY = cursorY
		if e.scrollBar.pos < 0 {
			e.scrollBar.pos = 0
		}
		if e.scrollBar.pos > e.height-e.scrollBar.height {
			e.scrollBar.pos = e.height - e.scrollBar.height
		}
	}

	maxPos := float64(e.height - e.scrollBar.height)
	pos := float64(e.scrollBar.pos)
	e.scrollOffset = -(pos / maxPos) * float64(e.totalHeight-e.height)

	return nil
}

func (e *Explorer) Draw(screen *ebiten.Image) {
	// Fill background.
	bg := ebiten.NewImage(screen.Bounds().Dx()-constant.MenuWidth, e.height)
	bg.Fill(color.Gray{Y: constant.ExplorerGrayY})
	bgOp := &ebiten.DrawImageOptions{}
	bgX := constant.MenuWidth
	bgY := 0
	bgOp.GeoM.Translate(float64(bgX), float64(bgY))
	defer screen.DrawImage(bg, bgOp)

	spritesPerRow := (bg.Bounds().Dx() - e.offsetX) / (e.size + e.offsetX)
	numOfSprites := len(e.sprites)
	if numOfSprites == 0 {
		return
	}

	// ScrollBar.
	if e.scrollBar.show {
		scrollBar := ebiten.NewImage(e.scrollBar.width, bg.Bounds().Dy())
		scrollBar.Fill(color.Gray{Y: constant.ScrollBarGrayY})
		scrollBarOp := &ebiten.DrawImageOptions{}
		scrollBarX := bg.Bounds().Dx() - e.scrollBar.width
		scrollBarY := 0
		scrollBarOp.GeoM.Translate(float64(scrollBarX), float64(scrollBarY))
		bg.DrawImage(scrollBar, scrollBarOp)
		// ScrollBar handle.
		handleWidth := e.scrollBar.width
		handle := ebiten.NewImage(handleWidth, e.scrollBar.height)
		handle.Fill(color.Gray{Y: constant.ScrollBarHandleGrayY})
		handleOp := &ebiten.DrawImageOptions{}
		handleX := scrollBarX
		handleOp.GeoM.Translate(float64(handleX), float64(e.scrollBar.pos))
		bg.DrawImage(handle, handleOp)
	}

	// Sprites.
	originFrame := ebiten.NewImage(e.size, e.size)
	originFrame.Fill(color.White)
	boxGray := ebiten.NewImage(e.size/5, e.size/5)
	boxGray.Fill(color.Gray{Y: constant.ExplorerGrayY})
	boxTransparent := ebiten.NewImage(e.size/5, e.size/5)
	boxTransparent.Fill(color.Transparent)

	// Transparent.
	{
		isGray := false
		for x := 0; x < 5; x++ {
			for y := 0; y < 5; y++ {
				box := ebiten.NewImage(e.size/5, e.size/5)
				boxOp := &ebiten.DrawImageOptions{}
				boxOp.GeoM.Translate(float64(x*box.Bounds().Dx()), float64(y*box.Bounds().Dy()))
				if isGray {
					originFrame.DrawImage(boxGray, boxOp)
				} else {
					originFrame.DrawImage(boxTransparent, boxOp)
				}
				originFrame.DrawImage(box, boxOp)
				isGray = !isGray
			}
		}
	}
	frameLine := ebiten.NewImage(e.size+2, e.size+2)
	frameLine.Fill(color.Black)

	for i, sprite := range e.sprites {
		frame := ebiten.NewImage(e.size, e.size)
		frame.DrawImage(originFrame, nil)
		frameOp := &ebiten.DrawImageOptions{}
		frameLineOp := &ebiten.DrawImageOptions{}
		row := i % spritesPerRow
		col := i / spritesPerRow
		x := e.offsetX*(row+1) + (e.size * row)
		y := e.offsetY*(col+1) + (e.size * col)
		frameOp.GeoM.Translate(float64(x), float64(y))
		frameLineOp.GeoM.Translate(float64(x-1), float64(y-1))
		// Draw Sprite.
		if !sprite.IsEmpty() {
			img := sprite.Image
			width := img.Bounds().Dx()
			height := img.Bounds().Dy()
			scale := 1.0
			if width < e.size && height < e.size {
				if width > height {
					scale = float64(e.size / width)
				} else {
					scale = float64(e.size / height)
				}
			}
			if width > e.size || height > e.size {
				if width > height {
					scale = float64(e.size) / float64(width)
				} else {
					scale = float64(e.size) / float64(height)
				}
			}
			spriteOp := &ebiten.DrawImageOptions{}
			spriteOp.GeoM.Scale(float64(scale), float64(scale))
			spriteOp.GeoM.Translate(-(float64(width)*scale)/2, -(float64(height)*scale)/2)
			spriteOp.GeoM.Translate(float64(e.size/2), float64(e.size/2))
			if e.scrollBar.show {
				frameLineOp.GeoM.Translate(0, e.scrollOffset)
				frameOp.GeoM.Translate(0, e.scrollOffset)
			}
			frame.DrawImage(img, spriteOp)
		}
		if i == e.cursolOn {
			frameLine := ebiten.NewImage(e.size+2, e.size+2)
			frameLine.Fill(color.RGBA{R: 0, G: 0, B: 255, A: 255})
			bg.DrawImage(frameLine, frameLineOp)
		} else {
			bg.DrawImage(frameLine, frameLineOp)
		}
		bg.DrawImage(frame, frameOp)
	}
}

func (e *Explorer) Layout(outsideWidth, outsideHeight int) {
	// Resize
	e.height = outsideHeight / 3
}
