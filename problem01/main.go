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
	var sum int64

	s := bufio.NewScanner(os.Stdin)
	for i := 1; s.Scan(); i++ {
		line := strings.TrimSpace(s.Text())
		mass, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Fatalf("could not read line %d (%q): %v", i, line, err)
		}
		sum += computeFuel(mass)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("total sum is", sum)
}

func computeFuel(mass int64) int64 {
	// take its mass, divide by three, round down, and subtract 2.
	return mass/3 - 2
}
