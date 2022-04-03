package main

import (
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"strings"
	"time"
)

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

var (
	h    = 10
	l    = 50
	step = time.Second / 10
)

func main() {
	snake := NewSnake(l/2, h/2, 1, RIGHT, l, h)
	g := NewSnakeGame(h, l, snake)

	inCh := make(chan []byte, 100)

	go func() {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
			return
		}

		defer term.Restore(int(os.Stdin.Fd()), oldState)

		b := make([]byte, 1)
		for {
			os.Stdin.Read(b)
			inCh <- b
		}
	}()

	stopCh := make(chan struct{})
	e := NewEngine(time.NewTicker(step), g, stopCh, inCh)
	go func() {
		e.Start()
	}()

	t := time.NewTicker(time.Second)
	i := 0
	go func() {
		for range t.C {
			i++
			if i == 100 {
				stopCh <- struct{}{}
				break
			}
		}
	}()

	<-stopCh
}

type Engine struct {
	t     *time.Ticker
	g     Game
	input chan []byte
	stop  chan struct{}
}

func NewEngine(t *time.Ticker, g Game, stop chan struct{}, input chan []byte) *Engine {
	return &Engine{
		t:     t,
		g:     g,
		stop:  stop,
		input: input,
	}
}

func (e *Engine) Start() {
	go func() {
		for {
			select {
			case <-e.t.C:
				// todo size-sensitive
				for i := 0; i < 11; i++ {
					fmt.Print("\r")
					if i != 10 {
						fmt.Print("\b")
					}
				}
				fmt.Print(e.g.String())

				e.g.Frame(e.input)
			case <-e.stop:
				return
			}
		}
	}()
}

func (e *Engine) Stop() {
	e.stop <- struct{}{}
}

func (e *Engine) redraw() {

}

type Game interface {
	fmt.Stringer
	Frame(input chan []byte)
}

type Snake struct {
	x int
	y int
	l int
	d Direction

	borderX int
	borderY int
}

func NewSnake(x, y, l int, d Direction, bx, by int) *Snake {
	return &Snake{
		x:       x,
		y:       y,
		l:       l,
		d:       d,
		borderX: bx,
		borderY: by,
	}
}

func (s *Snake) Move(in []byte) {
	switch string(in) {
	case "w":
		s.d = UP
	case "s":
		s.d = DOWN
	case "a":
		s.d = LEFT
	case "d":
		s.d = RIGHT
	default:

	}

	switch s.d {
	case UP:
		if s.y-1 == 0 {
			s.y = s.borderY - 2
		} else {
			s.y--
		}
	case DOWN:
		if s.y+1 == s.borderY-1 {
			s.y = 1
		} else {
			s.y++
		}
	case LEFT:
		if s.x-1 == 0 {
			s.x = s.borderX - 2
		} else {
			s.x--
		}
	case RIGHT:
		if s.x+1 == s.borderX-1 {
			s.x = 1
		} else {
			s.x++
		}
	}
}

type Field struct {
	field [][]string
	snake *Snake
}

func NewField(f [][]string, s *Snake) *Field {
	return &Field{
		field: f, snake: s,
	}
}

func (fld *Field) Redraw() {
	h := len(fld.field)
	w := len(fld.field[0])

	f := make([][]string, 0, h)
	for i := 0; i < h; i++ {
		ff := make([]string, 0, h)

		for j := 0; j < w; j++ {
			c := " "
			if i == 0 || i == h-1 || j == 0 || j == w-1 {
				c = "█"
			}

			if i == fld.snake.y && j == fld.snake.x {
				c = "*"
			}

			ff = append(ff, c)
		}

		f = append(f, ff)
	}

	fld.field = f
}

type SnakeGame struct {
	f *Field
}

func NewSnakeGame(h, w int, s *Snake) *SnakeGame {
	f := make([][]string, 0, h)

	for i := 0; i < h; i++ {
		ff := make([]string, 0, h)

		for j := 0; j < w; j++ {
			c := " "
			if i == 0 || i == h-1 || j == 0 || j == w-1 {
				c = "█"
			}

			if i == s.y && j == s.x {
				c = "*"
			}

			ff = append(ff, c)
		}

		f = append(f, ff)
	}

	return &SnakeGame{
		f: NewField(f, s),
	}
}

func (g *SnakeGame) String() string {
	var s strings.Builder
	for _, f := range g.f.field {
		for _, ff := range f {
			s.WriteString(ff)
		}
		s.WriteString("\r\n")
	}

	return s.String()
}

func (g *SnakeGame) Frame(input chan []byte) {
	select {
	case in := <-input:
		g.f.snake.Move(in)
	default:
		g.f.snake.Move([]byte("will not happen if not press button"))
	}
	g.f.Redraw()
}
