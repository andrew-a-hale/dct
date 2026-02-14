package art

import (
	"bytes"
	"log"

	"dct/cmd/utils"
)

var (
	NEWLINE []byte = []byte("\n")
	SPACE   []byte = []byte(" ")
)

var DCT []byte = []byte(
	`       __       __ 
  ____/ /_____ / /_
 / __  // ___// __/
/ /_/ // /__ / /_  
\__,_/ \___/ \__/  `,
)

var DUCKDB []byte = []byte(
	`       __              __        __ __  
  ____/ /__  __ _____ / /__ ____/ // /_ 
 / __  // / / // ___// //_// __  // __ \
/ /_/ // /_/ // /__ / ,<  / /_/ // /_/ /
\__,_/ \__,_/ \___//_/|_| \__,_//_.___/ `,
)

var CHARM []byte = []byte(
	`        __                           
  _____/ /_  ____ __________ ___     
 / ___/ __ \/ __ '/ ___/ __ '__ \    
/ /__/ / / / /_/ / /  / / / / / /    
\___/_/ /_/\__,_/_/  /_/ /_/ /_/_____
                              /_____/`,
)

const (
	LEFT = iota*2 - 1
	RIGHT
)

const (
	UP = iota*2 - 1
	DOWN
)

type Direction struct {
	X int
	Y int
}

type Position struct {
	Row int
	Col int
}

type Path []Position

type Size struct {
	Height int
	Width  int
}

type (
	Graphic struct {
		Art       [][][]byte
		Pos       Position
		Size      Size
		Direction Direction
	}
)

func makeGraphic(art []byte, row, col, windowWidth, windowHeight, dirX, dirY int) *Graphic {
	pixels := Pixels(art, windowWidth, windowHeight)
	pos := Position{row, col}
	size := Size{len(pixels), len(pixels[0])}

	// move graphic up and left if it doesn't fit
	if row+size.Height > windowHeight {
		pos.Row = windowHeight - size.Height
	}
	if col+size.Width > windowWidth {
		pos.Col = windowWidth - size.Width
	}

	dir := Direction{dirX, dirY}
	return &Graphic{pixels, pos, size, dir}
}

func (g Graphic) GetPos() (rowStart, colStart, rowEnd, colEnd int) {
	rowStart = g.Pos.Row
	rowEnd = rowStart + g.Size.Height
	colStart = g.Pos.Col
	colEnd = colStart + g.Size.Width
	return
}

func (g *Graphic) Update(scene *Scene) {
	gRowStart, gColStart, gRowEnd, gColEnd := g.GetPos()

	switch {
	case gRowStart < 0:
		fallthrough
	case gRowEnd > scene.Height:
		fallthrough
	case gColStart < 0:
		fallthrough
	case gColEnd > scene.Width:
		log.Fatal("out of bounds")
	}

	if gRowStart == 0 {
		g.Direction.Y = DOWN
	}

	if gRowEnd == scene.Height {
		g.Direction.Y = UP
	}

	if gColStart == 0 {
		g.Direction.X = RIGHT
	}

	if gColEnd == scene.Width {
		g.Direction.X = LEFT
	}

	g.Pos.Row += int(g.Direction.Y)
	g.Pos.Col += int(g.Direction.X)
}

func (g *Graphic) getPixel(row, col int) (char []byte, found bool) {
	rowStart, colStart, rowEnd, colEnd := g.GetPos()

	if row >= rowStart && row < rowEnd {
		if col >= colStart && col < colEnd {
			return g.Art[row-rowStart][col-colStart], true
		}
	}

	return SPACE, false
}

func Pixels(g []byte, width, height int) [][][]byte {
	lines := bytes.Split(g, NEWLINE)

	graphicWidth := len(lines[0])
	utils.Assert(graphicWidth < width && graphicWidth > 0, "graphic is bigger than display")

	graphicHeight := len(lines)
	utils.Assert(graphicHeight < height && graphicHeight > 0, "graphic is bigger than display")

	var cells [][][]byte
	for _, row := range lines {
		cells = append(cells, bytes.Split(row, nil))
	}

	return cells
}
