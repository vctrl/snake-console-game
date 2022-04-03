package main

import (
	"fmt"
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

func main() {
	snake := NewSnake(50/2, 10/2, 1, RIGHT)
	g := NewSnakeGame(10, 50, snake)

	stopCh := make(chan struct{})
	e := NewEngine(time.NewTicker(time.Second), g, stopCh)
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
	t    *time.Ticker
	g    Game
	stop chan struct{}
}

func NewEngine(t *time.Ticker, g Game, stop chan struct{}) *Engine {
	return &Engine{
		t:    t,
		g:    g,
		stop: stop,
	}
}

func (e *Engine) Start() {
	go func() {
		for {
			select {
			case <-e.t.C:
				// todo size-sensitive
				fmt.Print("\r\b\r\b\r\b\r\b\r\b\r\b\r\b\r\b\r\b\r\b\r\b")
				fmt.Print(e.g.String())
				e.g.Frame()
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
	Frame()
}

type Snake struct {
	x int
	y int
	l int
	d Direction
}

func NewSnake(x, y, l int, d Direction) *Snake {
	return &Snake{
		x: x,
		y: y,
		l: l,
		d: d,
	}
}

func (s *Snake) Move() {
	switch s.d {
	case UP:
		s.y++
	case DOWN:
		s.y--
	case LEFT:
		s.x--
	case RIGHT:
		s.x++
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
		s.WriteString("\n")
	}

	return s.String()
}

func (g *SnakeGame) Frame() {
	g.f.snake.Move()
	g.f.Redraw()
}
