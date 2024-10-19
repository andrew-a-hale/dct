package art

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type Scene struct {
	Graphic *Graphic
	Height  int
	Width   int
}

const FRAMERATE int = 15

func getGraphics() map[int]*[]byte {
	graphics := make(map[int]*[]byte, 3)
	graphics[0] = &DCT
	graphics[1] = &DUCKDB
	graphics[2] = &CHARM
	return graphics
}

func fpsToDuration() time.Duration {
	return time.Duration(1000/FRAMERATE) * time.Millisecond
}

var ArtCmd = &cobra.Command{
	Use:   "art",
	Short: "some dct art",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		graphics := getGraphics()
		scene := newScene()

		frame := 1
		var i int
		for {
			if frame%(FRAMERATE*2) == 0 {
				i = (i + 1) % len(graphics)
				scene.Graphic = makeGraphic(
					*graphics[i],
					scene.Graphic.Pos.Row,
					scene.Graphic.Pos.Col,
					scene.Width,
					scene.Height,
					int(scene.Graphic.Direction.X),
					int(scene.Graphic.Direction.Y),
				)
			}
			scene.Draw()
			scene.Update(scene.Graphic)
			time.Sleep(fpsToDuration())
			frame++
		}
	},
}

func newScene() Scene {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln("failed to get terminal size")
	}

	dct := makeGraphic(DCT, 3, 3, width, height, LEFT, UP)
	return Scene{dct, height, width}
}

func (s *Scene) Update(g *Graphic) error {
	g.Update(s)

	return nil
}

func (s *Scene) Draw() error {
	var char []byte
	var out string

	for i := 0; i < s.Height; i++ {
		for j := 0; j < s.Width; j++ {
			char, _ = s.Graphic.getPixel(i, j)
			out += fmt.Sprintf("%s", char)
		}
	}

	fmt.Print(out)

	return nil
}
