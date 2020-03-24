package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	// text, err := ioutil.ReadAll(os.Stdin)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	text := "3,0,4,0,99"

	nums := strings.Split(string(text), ",")
	program := make([]int, len(nums))
	for i, num := range nums {
		v, err := strconv.Atoi(strings.TrimSpace(num))
		if err != nil {
			log.Fatalf("could not parse number %q: %v", num, err)
		}
		program[i] = v
	}

	c := newComputer(program)
	for !c.done {
		fmt.Println(c)
		if err := c.next(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(c)
}

type computer struct {
	cells    []int
	nextInst int
	done     bool
}

func newComputer(program []int) *computer {
	cells := make([]int, len(program))
	copy(cells, program)
	return &computer{cells: cells}
}

func (c computer) String() string {
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

func (c *computer) read(pos int) int   { return c.cells[pos] }
func (c *computer) write(pos, val int) { c.cells[pos] = val }

type opCode int

const (
	opAdd    opCode = 1
	opMult   opCode = 2
	opInput  opCode = 3 // uses dest only
	opOutput opCode = 4 // uses src1 only
	opHalt   opCode = 99
)

type instruction struct {
	code             opCode
	src1, src2, dest int
}

func (i instruction) String() string {
	switch i.code {
	case opHalt:
		return "HALT"
	case opAdd:
		return fmt.Sprintf("ADD %2d = %2d + %2d", i.dest, i.src1, i.src2)
	case opMult:
		return fmt.Sprintf("MUL %2d = %2d + %2d", i.dest, i.src1, i.src2)
	case opInput:
		return fmt.Sprintf("INPUT %2d", i.src1)
	case opOutput:
		return fmt.Sprintf("OUTPUT %2d", i.dest)
	}
	return "UNKOWN"
}

func (c *computer) nextInstruction() (*instruction, error) {
	op := opCode(c.cells[c.nextInst])
	switch op {
	case opHalt:
		return &instruction{code: op}, nil
	case opAdd, opMult:
		idx := c.nextInst
		c.nextInst += 4
		return &instruction{
			code: op,
			src1: c.cells[idx+1],
			src2: c.cells[idx+2],
			dest: c.cells[idx+3],
		}, nil
	case opInput, opOutput:
		idx := c.nextInst
		c.nextInst += 2

		ins := &instruction{code: op}
		if op == opInput {
			ins.dest = c.cells[idx+1]
		} else {
			ins.src1 = c.cells[idx+1]
		}
		return ins, nil
	default:
		return nil, fmt.Errorf("unknown op code %d", op)
	}
}

func (c *computer) next() error {
	inst, err := c.nextInstruction()
	if err != nil {
		return err
	}

	switch inst.code {
	case opAdd:
		a := c.read(inst.src1)
		b := c.read(inst.src2)
		c.write(inst.dest, a+b)
	case opMult:
		a := c.read(inst.src1)
		b := c.read(inst.src2)
		c.write(inst.dest, a*b)
	case opInput:
		// TODO: this won't always be one
		c.write(inst.dest, 1)
	case opOutput:
		fmt.Println("OUTPUT:", c.read(inst.src1))
	case opHalt:
		c.done = true
	}
	return nil
}
