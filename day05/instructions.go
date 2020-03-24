package main

import "fmt"

type OpCode int

const (
	OpAdd    OpCode = 1
	OpMult   OpCode = 2
	OpInput  OpCode = 3 // uses dest only
	OpOutput OpCode = 4 // uses src1 only
	OpHalt   OpCode = 99
)

func NewInstruction(op OpCode) (Instruction, error) {
	switch op {
	case OpAdd:
		return new(addInstruction), nil
	case OpMult:
		return new(multInstruction), nil
	case OpInput:
		return new(inputInstruction), nil
	case OpOutput:
		return new(outputInstruction), nil
	case OpHalt:
		return new(haltInstruction), nil
	default:
		return nil, fmt.Errorf("unknown op code %d", op)
	}
}

type Instruction interface {
	Op() OpCode
	Parse(c *Computer) error
	Run(c *Computer)
	String() string
}

type addInstruction struct {
	src1 int
	src2 int
	dest int
}

func (i *addInstruction) Op() OpCode { return OpAdd }

func (i *addInstruction) Parse(c *Computer) error {
	idx := c.nextInst
	c.nextInst += 4
	i.src1 = c.cells[idx+1]
	i.src2 = c.cells[idx+2]
	i.dest = c.cells[idx+3]
	return nil
}

func (i *addInstruction) String() string {
	return fmt.Sprintf("ADD %2d = %2d + %2d", i.dest, i.src1, i.src2)
}

func (i *addInstruction) Run(c *Computer) {
	a := c.read(i.src1)
	b := c.read(i.src2)
	c.write(i.dest, a+b)
}

type multInstruction struct {
	src1 int
	src2 int
	dest int
}

func (i *multInstruction) Op() OpCode { return OpMult }

func (i *multInstruction) Parse(c *Computer) error {
	idx := c.nextInst
	c.nextInst += 4
	i.src1 = c.cells[idx+1]
	i.src2 = c.cells[idx+2]
	i.dest = c.cells[idx+3]
	return nil
}

func (i *multInstruction) String() string {
	return fmt.Sprintf("MUL %2d = %2d * %2d", i.dest, i.src1, i.src2)
}

func (i *multInstruction) Run(c *Computer) {
	a := c.read(i.src1)
	b := c.read(i.src2)
	c.write(i.dest, a*b)
}

type outputInstruction struct {
	src int
}

func (i *outputInstruction) Op() OpCode { return OpOutput }

func (i *outputInstruction) Parse(c *Computer) error {
	idx := c.nextInst
	c.nextInst += 2
	i.src = c.cells[idx+1]
	return nil
}

func (i *outputInstruction) String() string {
	return fmt.Sprintf("OUTPUT %2d", i.src)
}

func (i *outputInstruction) Run(c *Computer) {
	// TODO: this won't always be one
	fmt.Println("OUTPUT:", c.read(i.src))
}

type inputInstruction struct {
	dest int
}

func (i *inputInstruction) Op() OpCode { return OpInput }

func (i *inputInstruction) Parse(c *Computer) error {
	idx := c.nextInst
	c.nextInst += 2
	i.dest = c.cells[idx+1]
	return nil
}

func (i *inputInstruction) String() string {
	return fmt.Sprintf("INPUT %2d", i.dest)
}

func (i *inputInstruction) Run(c *Computer) {
	// TODO: this won't always be one
	c.write(i.dest, 1)
}

type haltInstruction struct{}

func (i *haltInstruction) Op() OpCode              { return OpHalt }
func (i *haltInstruction) Parse(c *Computer) error { return nil }
func (i *haltInstruction) String() string          { return "HALT" }
func (i *haltInstruction) Run(c *Computer)         { c.done = true }
