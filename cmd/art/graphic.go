package art

import (
	"bytes"
	"dct/cmd/utils"

	"golang.org/x/exp/rand"
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

type Graphic struct {
	Id     string
	Art    [][][]byte
	Pos    [2]int // row, col
	Size   [2]int // width, height
	Dir    [2]int // vertical, horizontal
	Speed  int
	Render bool
}

func makeGraphic(id string, art []byte, width, height, row, col, speed int) *Graphic {
	pixels := Pixels(art, width, height)
	pos := [2]int{row, col}
	size := [2]int{len(pixels), len(pixels[0])}
	dir := [2]int{rand.Intn(2)*2 - 1, rand.Intn(2)*2 - 1}
	return &Graphic{id, pixels, pos, size, dir, speed, false}
}

func (g *Graphic) boundCheck(scene *Scene) {
	gRowStart := g.Pos[0] + g.Speed*g.Dir[0]
	gColStart := g.Pos[1] + g.Speed*g.Dir[1]
	gRowEnd := g.Pos[0] + g.Size[0] + g.Speed*g.Dir[0]
	gColEnd := g.Pos[1] + g.Size[1] + g.Speed*g.Dir[1]

	if gRowStart < 0 || gRowEnd > scene.Height {
		g.Dir[0] = g.Dir[0] * -1
	}

	if gColStart < 0 || gColEnd > scene.Width {
		g.Dir[1] = g.Dir[1] * -1
	}
}

func (g *Graphic) checkCollision(o *Graphic) bool {
	if g.Pos[1]+g.Size[1]+g.Speed*g.Dir[1] > o.Pos[1] &&
		g.Pos[1]+g.Speed*g.Dir[1] < o.Pos[1]+o.Size[1] &&
		g.Pos[0]+g.Size[0] > o.Pos[0] &&
		g.Pos[0] < o.Pos[0]+o.Size[0] {
		g.Dir[1] = g.Dir[1] * -1
		o.Dir[1] = o.Dir[1] * -1
	}

	if g.Pos[1]+g.Size[1] > o.Pos[1] &&
		g.Pos[1] < o.Pos[1]+o.Size[1] &&
		g.Pos[0]+g.Size[0]+g.Speed*g.Dir[0] > o.Pos[0] &&
		g.Pos[0]+g.Speed*g.Dir[0] < o.Pos[0]+o.Size[0] {
		g.Dir[0] = g.Dir[0] * -1
		o.Dir[0] = o.Dir[0] * -1
	}

	return false
}

func (g *Graphic) getPixel(row, col int) (char []byte, found bool) {
	if row >= g.Pos[0] && row < g.Pos[0]+g.Size[0] {
		if col >= g.Pos[1] && col < g.Pos[1]+g.Size[1] {
			char = g.Art[row-g.Pos[0]][col-g.Pos[1]]
			return char, true
		}
	}

	return SPACE, false
}

func Pixels(g []byte, width, height int) [][][]byte {
	lines := bytes.Split(g, NEWLINE)

	graphicWidth := len(lines[0])
	utils.Assert(graphicWidth < width && graphicWidth > 0, "graphicWidth out of bounds")

	graphicHeight := len(lines)
	utils.Assert(graphicHeight < height && graphicHeight > 0, "graphicHeight out of bounds")

	var cells [][][]byte
	for _, row := range lines {
		cells = append(cells, bytes.Split(row, nil))
	}

	return cells
}
