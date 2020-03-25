package main

import (
	"fmt"
	"log"
)

type opCode int

const (
	opAdd         opCode = 1
	opMult        opCode = 2
	opInput       opCode = 3
	opOutput      opCode = 4
	opJumpIfTrue  opCode = 5
	opJumpIfFalse opCode = 6
	opLessThan    opCode = 7
	opEquals      opCode = 8
	opHalt        opCode = 99
)

type Instruction interface {
	Parse(c *Computer)
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
	case opJumpIfTrue, opJumpIfFalse:
		return &condJumpInstruction{jumpOn: op == opJumpIfTrue}, nil
	case opLessThan:
		return new(lessThanInstruction), nil
	case opEquals:
		return new(equalsInstruction), nil
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

func (i *outputInstruction) String() string { return fmt.Sprintf("OUTPUT %v", i.arg) }

func (i *outputInstruction) Run(c *Computer) { fmt.Println("OUTPUT:", i.arg.read(c)) }

type inputInstruction struct{ unaryOpInstruction }

func (i *inputInstruction) String() string { return fmt.Sprintf("INPUT %v", i.arg) }

func (i *inputInstruction) Run(c *Computer) {
	val := 0
	_, err := fmt.Scanf("%d", &val)
	if err != nil {
		log.Fatalf("could not read from input: %v", err)
	}
	i.arg.write(c, val)
}

type condJumpInstruction struct {
	jumpOn bool
	cond   parameter
	target parameter
}

func (i *condJumpInstruction) Parse(c *Computer) {
	params := parseParameters(c, 2)
	i.cond = params[0]
	i.target = params[1]
}

func (i *condJumpInstruction) String() string {
	return fmt.Sprintf("JumpIf(%v) %v %v", i.jumpOn, i.cond, i.target)
}

func (i *condJumpInstruction) Run(c *Computer) {
	if (i.cond.read(c) == 0) == i.jumpOn {
		return
	}
	c.nextInst = i.target.read(c)
}

type lessThanInstruction struct{ binaryOpInstruction }

func (i *lessThanInstruction) String() string {
	return fmt.Sprintf("LessThan: %v = %v < %v", i.dest, i.src1, i.src2)
}

func (i *lessThanInstruction) Run(c *Computer) {
	if i.src1.read(c) < i.src2.read(c) {
		i.dest.write(c, 1)
	} else {
		i.dest.write(c, 0)
	}
}

type equalsInstruction struct{ binaryOpInstruction }

func (i *equalsInstruction) String() string {
	return fmt.Sprintf("Equals: %v = %v == %v", i.dest, i.src1, i.src2)
}

func (i *equalsInstruction) Run(c *Computer) {
	if i.src1.read(c) == i.src2.read(c) {
		i.dest.write(c, 1)
	} else {
		i.dest.write(c, 0)
	}
}

type haltInstruction struct{}

func (i *haltInstruction) Parse(c *Computer) {}
func (i *haltInstruction) String() string    { return "HALT" }
func (i *haltInstruction) Run(c *Computer)   { c.done = true }

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

func (i *unaryOpInstruction) Parse(c *Computer) {
	i.arg = parseParameters(c, 1)[0]
}

type binaryOpInstruction struct {
	src1, src2, dest parameter
}

func (i *binaryOpInstruction) Parse(c *Computer) {
	params := parseParameters(c, 3)
	i.src1 = params[0]
	i.src2 = params[1]
	i.dest = params[2]
}

func parseParameters(c *Computer, n int) []parameter {
	idx := c.nextInst
	c.nextInst += n + 1

	modes := c.cells[idx] / 100

	var params []parameter
	for i := 1; i <= n; i++ {
		params = append(params, parameter{
			c.cells[idx+i],
			modes%10 == 1,
		})
		modes = modes / 10
	}

	return params
}
