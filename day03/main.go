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

	res, dist, ok := mergeBoards(cable1, cable2)
	if !ok {
		log.Fatalf("could not find any intersections")
	}
	fmt.Printf("intersection with lowest cable distance found at %v with distance %d\n", res, dist)
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

type board map[position]int

func parseCable(s *bufio.Scanner) (board, error) {
	if !s.Scan() {
		return nil, fmt.Errorf("no more lines in input")
	}
	line := s.Text()
	moves := strings.Split(line, ",")

	b := make(board)
	var p position

	dist := 0
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
			dist++
			if b[p] != 0 {
				continue
			}
			b[p] = dist
		}
	}

	return b, nil
}

func mergeBoards(b1, b2 board) (position, int, bool) {
	var minInt position
	var minDist int

	for p := range b1 {
		if p.isZero() || b1[p] == 0 || b2[p] == 0 {
			continue
		}
		fmt.Printf("found intersection at %v\n", p)
		if dist := b1[p] + b2[p]; minInt.isZero() || dist < minDist {
			minInt = p
			minDist = dist
		}
	}
	return minInt, minDist, !minInt.isZero()
}
