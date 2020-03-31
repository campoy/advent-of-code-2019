package intcode

import (
	"bytes"
	"fmt"
)

type Computer struct {
	cells    []int
	nextInst int
	done     bool
	stdin    chan int
	stdout   chan int
}

func NewComputer(program []int, stdin, stdout chan int) *Computer {
	cells := make([]int, len(program))
	copy(cells, program)
	return &Computer{cells: cells, stdin: stdin, stdout: stdout}
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

func (c *Computer) Run() error {
	for !c.done {
		if err := c.next(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Computer) next() error {
	ins, err := newInstruction(c.cells[c.nextInst])
	if err != nil {
		return err
	}
	ins.parse(c)
	ins.run(c)
	return nil
}

func (c *Computer) read(pos int) int { return c.cells[pos] }

func (c *Computer) write(pos, val int) { c.cells[pos] = val }

func (c *Computer) Stdin() chan<- int  { return c.stdin }
func (c *Computer) Stdout() <-chan int { return c.stdout }
