package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type object struct {
	name      string
	orbits    *object
	orbitedBy []*object
}

func main() {
	objects := make(map[string]*object)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		pieces := strings.Split(s.Text(), ")")

		center := pieces[0]
		if objects[center] == nil {
			objects[center] = &object{name: center}
		}

		moon := pieces[1]
		if objects[moon] == nil {
			objects[moon] = &object{name: moon}
		}

		objects[center].orbitedBy = append(objects[center].orbitedBy, objects[moon])
		objects[moon].orbits = objects[center]
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	you := objects["YOU"]
	youOrbits := orbitsToCoM(you)

	san := objects["SAN"]
	sanOrbits := orbitsToCoM(san)

	for youOrbits[len(youOrbits)-1] == sanOrbits[len(sanOrbits)-1] {
		youOrbits = youOrbits[:len(youOrbits)-1]
		sanOrbits = sanOrbits[:len(sanOrbits)-1]
	}

	for _, o := range youOrbits {
		fmt.Printf("%s ", o.name)
	}
	fmt.Println()

	for _, o := range sanOrbits {
		fmt.Printf("%s ", o.name)
	}
	fmt.Println()

	fmt.Println(len(youOrbits) - 1 + len(sanOrbits) - 1)
}

func orbitsToCoM(obj *object) []*object {
	var path []*object
	for obj != nil {
		path = append(path, obj)
		obj = obj.orbits
	}
	return path
}
