package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/campoy/advent-of-code-2019/day07/intcode"
)

func main() {
	amplifiers := flag.Int("n", 5, "number of amplifiers")
	path := flag.String("i", "input.txt", "input program")
	flag.Parse()

	bs, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	var program []int
	for _, text := range bytes.Split(bs, []byte{','}) {
		code, err := strconv.Atoi(string(text))
		if err != nil {
			log.Fatalf("could not parse number %s: %v", text, err)
		}
		program = append(program, code)
	}

	perms := permutations([]int{0, 1, 2, 3, 4})

	maxResult := 0
	maxSettings := make([]int, *amplifiers)
	for _, settings := range perms {
		result, err := runWithSettings(program, *amplifiers, settings)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(settings, result)
		if result > maxResult {
			maxResult = result
			copy(maxSettings, settings)
		}
	}

	fmt.Printf("max result was %d with settings %v\n", maxResult, maxSettings)
}

func runWithSettings(program []int, amplifiers int, settings []int) (int, error) {
	input := make(chan int, 1)
	output := make(chan int)
	errc := make(chan error, amplifiers)

	ch := input
	for i := 0; i < amplifiers; i++ {
		ch <- settings[i]
		out := output
		if i < amplifiers-1 {
			out = make(chan int, 1)
		}

		c := intcode.NewComputer(program, ch, out)

		go func() { errc <- c.Run() }()
		ch = out
	}

	input <- 0
	for {
		select {
		case out := <-output:
			return out, nil
		case err := <-errc:
			if err != nil {
				return 0, err
			}
		}
	}
}

func permutations(values []int) [][]int {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return [][]int{values}
	}

	var perms [][]int
	for i, v := range values {
		rest := make([]int, len(values)-1)
		copy(rest[:i], values[:i])
		copy(rest[i:], values[i+1:])
		for _, perm := range permutations(rest) {
			perm = append([]int{v}, perm...)
			perms = append(perms, perm)
		}
	}

	return perms
}
