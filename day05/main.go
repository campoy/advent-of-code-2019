package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	nums := strings.Split(string(text), ",")
	program := make([]int, len(nums))
	for i, num := range nums {
		v, err := strconv.Atoi(strings.TrimSpace(num))
		if err != nil {
			log.Fatalf("could not parse number %q: %v", num, err)
		}
		program[i] = v
	}

	c := NewComputer(program)
	for !c.done {
		if err := c.next(); err != nil {
			log.Fatal(err)
		}
	}
}

type Computer struct {
	cells    []int
	nextInst int
	done     bool
}

func NewComputer(program []int) *Computer {
	cells := make([]int, len(program))
	copy(cells, program)
	return &Computer{cells: cells}
}

func (c *Computer) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "insptr:\t%d\n", c.nextInst)
	for i := 0; i < c.nextInst; i++ {
		fmt.Fprintf(w, "%5d ", c.cells[i])
		if i%10 == 9 {
			fmt.Fprintln(w)
		}
	}
	fmt.Fprintf(w, "|")
	for i := c.nextInst; i < len(c.cells); i++ {
		fmt.Fprintf(w, "%5d ", c.cells[i])
		if i%10 == 9 {
			fmt.Fprintln(w)
		}
	}
	fmt.Fprintln(w)

	return w.String()
}

func (c *Computer) next() error {
	ins, err := NewInstruction(c.cells[c.nextInst])
	if err != nil {
		return err
	}
	ins.Parse(c)
	ins.Run(c)
	return nil
}

func (c *Computer) read(pos int) int { return c.cells[pos] }

func (c *Computer) write(pos, val int) { c.cells[pos] = val }
