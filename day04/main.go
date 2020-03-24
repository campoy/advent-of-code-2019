package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "needs two numerical arguments, got %d\n", len(args))
		os.Exit(1)
	}

	from, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "1st argument could not be parsed as number: %v\n", err)
		os.Exit(1)
	}

	to, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "2nd argument could not be parsed as number: %v\n", err)
		os.Exit(1)
	}

	count := 0
	last := newPassword(to)
	for pass := newPassword(from); !pass.equals(last); pass.next() {
		if pass.isValid() {
			count++
		}
	}
	fmt.Println(count)
}

type password struct {
	digits [6]int
}

func newPassword(val int) *password {
	var p password
	for i := len(p.digits) - 1; i >= 0; i-- {
		p.digits[i] = val % 10
		val = val / 10
	}
	return &p
}

func (p *password) String() string {
	d := p.digits
	return fmt.Sprintf("%d%d%d%d%d%d", d[0], d[1], d[2], d[3], d[4], d[5])
}

func (p *password) equals(q *password) bool {
	for i := range p.digits {
		if p.digits[i] != q.digits[i] {
			return false
		}
	}
	return true
}

func (p *password) isValid() bool {
	repeated := false
	for i := 1; i < len(p.digits); i++ {
		if p.digits[i] < p.digits[i-1] {
			return false
		}
		if p.digits[i] == p.digits[i-1] {
			repeated = true
		}
	}
	return repeated
}

func (p *password) next() {
	p.digits[len(p.digits)-1]++

	for i := len(p.digits) - 1; i > 0; i-- {
		if p.digits[i] > 9 {
			p.digits[i] = 0
			p.digits[i-1]++
		}
	}

	if p.digits[0] > 9 {
		log.Fatal("password has too many digits")
	}
}
