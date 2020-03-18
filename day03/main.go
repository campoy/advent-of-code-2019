package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	cable1, err := parseCable(s)
	if err != nil {
		log.Fatalf("could not parse first cable: %v", err)
	}

	cable2, err := parseCable(s)
	if err != nil {
		log.Fatalf("could not parse first cable: %v", err)
	}

	res, ok := mergeBoards(cable1, cable2)
	if !ok {
		log.Fatalf("could not find any intersections")
	}
	fmt.Printf("closest intersection to origin is %v with distance %d\n", res, res.distToZero())
}

type position struct{ x, y int }

func (p *position) move(v position) {
	p.x += v.x
	p.y += v.y
}

func (p *position) isZero() bool { return p.x == 0 && p.y == 0 }

func (p *position) distToZero() int {
	dx := p.x
	if dx < 0 {
		dx = -dx
	}
	dy := p.y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

type board map[position]bool

func parseCable(s *bufio.Scanner) (board, error) {
	if !s.Scan() {
		return nil, fmt.Errorf("no more lines in input")
	}
	line := s.Text()
	moves := strings.Split(line, ",")

	b := make(board)
	var p position

	b[p] = true

	for i, move := range moves {
		var inc position
		switch dir := move[0]; dir {
		case 'R':
			inc = position{1, 0}
		case 'L':
			inc = position{-1, 0}
		case 'U':
			inc = position{0, 1}
		case 'D':
			inc = position{0, -1}
		}

		steps, err := strconv.Atoi(move[1:])
		if err != nil {
			return nil, fmt.Errorf("could not parse number of steps on move %d: %w", i, err)
		}

		for i := 0; i < steps; i++ {
			p.move(inc)
			b[p] = true
		}
	}

	return b, nil
}

func mergeBoards(b1, b2 board) (position, bool) {
	var minInt position
	for p, ok := range b1 {
		if p.isZero() || !ok || !b2[p] {
			continue
		}
		fmt.Printf("found intersection at %v\n", p)
		if minInt.isZero() || p.distToZero() < minInt.distToZero() {
			minInt = p
		}
	}
	return minInt, !minInt.isZero()
}
