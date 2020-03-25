package main

import (
	"fmt"
	"log"
)

type opCode int

const (
	opAdd    opCode = 1
	opMult   opCode = 2
	opInput  opCode = 3
	opOutput opCode = 4
	opHalt   opCode = 99
)

type Instruction interface {
	Parse(c *Computer) error
	Run(c *Computer)
	String() string
}

func NewInstruction(val int) (Instruction, error) {
	op := opCode(val % 100)
	switch op {
	case opAdd:
		return new(addInstruction), nil
	case opMult:
		return new(multInstruction), nil
	case opInput:
		return new(inputInstruction), nil
	case opOutput:
		return new(outputInstruction), nil
	case opHalt:
		return new(haltInstruction), nil
	default:
		return nil, fmt.Errorf("unknown op code %d", op)
	}
}

type addInstruction struct{ binaryOpInstruction }

func (i *addInstruction) Run(c *Computer) { i.dest.write(c, i.src1.read(c)+i.src2.read(c)) }

func (i *addInstruction) String() string {
	return fmt.Sprintf("ADD %v = %v + %v", i.dest, i.src1, i.src2)
}

type multInstruction struct{ binaryOpInstruction }

func (i *multInstruction) Run(c *Computer) { i.dest.write(c, i.src1.read(c)*i.src2.read(c)) }

func (i *multInstruction) String() string {
	return fmt.Sprintf("MUL %v = %v * %v", i.dest, i.src1, i.src2)
}

type outputInstruction struct{ unaryOpInstruction }

func (i *outputInstruction) String() string  { return fmt.Sprintf("OUTPUT %v", i.arg) }
func (i *outputInstruction) Run(c *Computer) { fmt.Println("OUTPUT:", i.arg.read(c)) }

type inputInstruction struct{ unaryOpInstruction }

func (i *inputInstruction) String() string  { return fmt.Sprintf("INPUT %v", i.arg) }
func (i *inputInstruction) Run(c *Computer) { i.arg.write(c, 1) }

type haltInstruction struct{}

func (i *haltInstruction) Parse(c *Computer) error { return nil }
func (i *haltInstruction) String() string          { return "HALT" }
func (i *haltInstruction) Run(c *Computer)         { c.done = true }

// utility instructions

type parameter struct {
	value     int
	immediate bool
}

func (p parameter) String() string {
	if p.immediate {
		return fmt.Sprint(p.value)
	}
	return fmt.Sprintf("@%d", p.value)
}

func (p parameter) read(c *Computer) int {
	if p.immediate {
		return p.value
	}
	return c.read(p.value)
}

func (p parameter) write(c *Computer, val int) {
	if p.immediate {
		log.Fatal("wrote into an immediate parameter")
	}
	c.write(p.value, val)
}

type unaryOpInstruction struct{ arg parameter }

func (i *unaryOpInstruction) Parse(c *Computer) error {
	idx := c.nextInst
	c.nextInst += 2
	immediate := (c.cells[idx] / 100) == 1
	i.arg = parameter{c.cells[idx+1], immediate}
	return nil
}

type binaryOpInstruction struct {
	src1, src2, dest parameter
}

func (i *binaryOpInstruction) Parse(c *Computer) error {
	idx := c.nextInst
	c.nextInst += 4

	op := c.cells[idx]
	i.src1 = parameter{c.cells[idx+1], (op/100)%10 == 1}
	i.src2 = parameter{c.cells[idx+2], (op/1000)%10 == 1}
	i.dest = parameter{c.cells[idx+3], op/10000 == 1}
	return nil
}
