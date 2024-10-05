package art

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const DCT string = `       __       __ 
  ____/ /_____ / /_
 / __  // ___// __/
/ /_/ // /__ / /_  
\__,_/ \___/ \__/  `

const DUCKDB string = `       __              __        __ __  
  ____/ /__  __ _____ / /__ ____/ // /_ 
 / __  // / / // ___// //_// __  // __ \
/ /_/ // /_/ // /__ / ,<  / /_/ // /_/ /
\__,_/ \__,_/ \___//_/|_| \__,_//_.___/ `

type Graphic struct {
	art   string
	pos   [2]int
	size  [2]int
	dir   [2]int
	speed int
}

type Scene struct {
	Graphics []Graphic
	Height   int
	Width    int
}

const FRAMERATE int = 60

var (
	height int
	width  int
)

var ArtCmd = &cobra.Command{
	Use:   "art",
	Short: "some dct art",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		scene := newScene()
		scene.Draw()
	},
}

func (g *Graphic) boundCheck(width, height int) bool {
	_, _ = width, height
	return false
}

func (g *Graphic) collide(o *Graphic) bool {
	_ = o
	return false
}

func makeDct(width, height int) Graphic {
	_, _ = width, height
	lines := strings.Split(DCT, "\n")
	return Graphic{DCT, [2]int{100, 100}, [2]int{len(lines[0]), len(lines)}, [2]int{1, 1}, 10}
}

func makeDuckDb(width, height int) Graphic {
	_, _ = width, height
	lines := strings.Split(DUCKDB, "\n")
	return Graphic{DUCKDB, [2]int{200, 200}, [2]int{len(lines[0]), len(lines)}, [2]int{1, -1}, 10}
}

func newScene() Scene {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln("failed to get terminal size")
	}
	dct := makeDct(width, height)
	duckdb := makeDuckDb(width, height)

	var graphics []Graphic
	graphics = append(graphics, dct, duckdb)
	return Scene{graphics, height, width}
}

func (s *Scene) Update() error {
	return nil
}

func (s *Scene) Draw() error {
	for _, g := range s.Graphics {
		fmt.Printf("%v\n", g)
	}

	return nil
}
