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
	cells := make([]int, len(nums))
	for i, num := range nums {
		cells[i], err = strconv.Atoi(strings.TrimSpace(num))
		if err != nil {
			log.Fatalf("could not parse number %q: %v", num, err)
		}
	}

	c := computer{cells: cells}
	for !c.done {
		fmt.Println(c)
		if err := c.next(); err != nil {
			log.Fatal(err)
		}
	}
}

type computer struct {
	cells    []int
	nextInst int
	done     bool
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

func (c *computer) read(pos int) int {
	fmt.Printf("r(%d): %d\n", pos, c.cells[pos])
	return c.cells[pos]
}

func (c *computer) write(pos, val int) {
	fmt.Printf("w(%d, %d)\n", pos, val)
	c.cells[pos] = val
}

type opCode int

const (
	opAdd  opCode = 1
	opMult opCode = 2
	opHalt opCode = 99
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
	}
	return "UNKOWN"
}

func (c *computer) nextInstruction() (*instruction, error) {
	op := opCode(c.cells[c.nextInst])
	switch op {
	case opHalt:
		c.nextInst++
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
	default:
		return nil, fmt.Errorf("unknown op code %d", op)
	}
}

func (c *computer) next() error {
	inst, err := c.nextInstruction()
	if err != nil {
		return err
	}

	fmt.Println("executing instruction:", inst)

	switch inst.code {
	case opAdd:
		a := c.read(inst.src1)
		b := c.read(inst.src2)
		c.write(inst.dest, a+b)
	case opMult:
		a := c.read(inst.src1)
		b := c.read(inst.src2)
		c.write(inst.dest, a*b)
	case opHalt:
		c.done = true
	}
	return nil
}
