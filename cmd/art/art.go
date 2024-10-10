package art

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/rand"
	"golang.org/x/term"
)

type Scene struct {
	Graphics []*Graphic
	Height   int
	Width    int
}

const FRAMERATE int = 20

func fpsToDuration() time.Duration {
	return time.Duration(1000/FRAMERATE) * time.Millisecond
}

var ArtCmd = &cobra.Command{
	Use:   "art",
	Short: "some dct art",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		scene := newScene()

		var frame int
		var i int
		for {
			if frame%200 == 0 {
				scene.Graphics[i].Render = false
				rand.Seed(uint64(time.Now().UnixNano()))
				i = rand.Intn(len(scene.Graphics))
				scene.Graphics[i].Render = true
			}
			scene.Draw()
			scene.UpdateSingle(scene.Graphics[i])
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

	var display [][][]byte
	for i := 0; i < height; i++ {
		var row [][]byte
		for j := 0; j < width; j++ {
			row = append(row, []byte(" "))
		}
		display = append(display, row)
	}

	dct := makeGraphic("dct", DCT, width, height, 10, 60, 1)
	duckdb := makeGraphic("duckdb", DUCKDB, width, height, 1, 20, 1)
	charm := makeGraphic("charm", CHARM, width, height, 20, 20, 1)

	var graphics []*Graphic
	graphics = append(graphics, dct, duckdb, charm)
	return Scene{graphics, height, width}
}

func (s *Scene) UpdateSingle(g *Graphic) error {
	g.boundCheck(s)
	g.Pos[0] = g.Pos[0] + g.Speed*g.Dir[0]
	g.Pos[1] = g.Pos[1] + g.Speed*g.Dir[1]

	return nil
}

func (s *Scene) UpdateMultiple() error {
	for i := 0; i < len(s.Graphics); i++ {
		g := s.Graphics[i]
		for j := i + 1; j < len(s.Graphics); j++ {
			g.checkCollision(s.Graphics[j])
		}
		g.boundCheck(s)
		g.Pos[0] = g.Pos[0] + g.Speed*g.Dir[0]
		g.Pos[1] = g.Pos[1] + g.Speed*g.Dir[1]
	}

	return nil
}

func (s *Scene) Draw() error {
	var char []byte
	var found bool
	var out string

	for i := 0; i < s.Height; i++ {
		for j := 0; j < s.Width; j++ {
			for _, g := range s.Graphics {
				if g.Render {
					char, found = g.getPixel(i, j)
					if found {
						break
					}
				}
			}
			out = out + fmt.Sprintf("%s", char)
		}
	}

	fmt.Print(out)

	return nil
}
