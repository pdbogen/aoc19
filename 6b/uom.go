package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Body string
type Orbit struct {
	Satellites map[Body]*Orbit
}

// Checksum returns the number of direct and indirect orbits are attached to
// the given orbit. A leaf has zero orbits.
func Checksum(orbit *Orbit) int {
	return checksumN(orbit, 0)
}

func checksumN(orbit *Orbit, n int) int {
	if len(orbit.Satellites) == 0 {
		return n
	}

	sum := n
	for _, satellite := range orbit.Satellites {
		sck := checksumN(satellite, n+1)
		sum += sck
	}
	return sum
}

// ParseMap parses an unordered list of `body)satellite` pairs
func ParseMap(in io.Reader) (*Orbit, error) {
	// maps a body to a list of satellites
	satellites := map[string][]string{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		txt := scanner.Text()
		orbit := strings.Split(txt, ")")
		if len(orbit) != 2 {
			return nil, fmt.Errorf("malformed orbit description %s", txt)
		}
		satellites[orbit[0]] = append(satellites[orbit[0]], orbit[1])
	}
	return toOrbits(satellites, Body("COM")), nil
}

func toOrbits(satellites map[string][]string, body Body) (*Orbit) {
	ret := &Orbit{map[Body]*Orbit{}}
	for _, sat := range satellites[string(body)] {
		ret.Satellites[Body(sat)] = toOrbits(satellites, Body(sat))
	}
	return ret
}

func ComPath(com *Orbit, target Body) []Body {
	if len(com.Satellites) == 0 {
		return nil
	}
	if _, ok := com.Satellites[target]; ok {
		return []Body{}
	}
	for satName, sat := range com.Satellites {
		res := ComPath(sat, target)
		if res == nil {
			continue
		}
		return append([]Body{satName}, res...)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("%s <uom data file>", os.Args[0])
	}

	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("reading %q: %v", os.Args[1], err)
	}

	comOrbit, err := ParseMap(data)
	if err != nil {
		log.Fatalf("parse uom data from %q: %v", os.Args[1], err)
	}

	log.Printf("UOM data checksum: %d", Checksum(comOrbit))

	santaPath := ComPath(comOrbit, "SAN")
	ourPath := ComPath(comOrbit, "YOU")
	log.Printf("Santa Path: %v", santaPath)
	log.Printf("Our Path: %v", ourPath)

	for {
		if santaPath[1] != ourPath[1] {
			break
		}
		santaPath = santaPath[1:]
		ourPath = ourPath[1:]
	}
	log.Printf("Santa Path: %v", santaPath)
	log.Printf("Our Path: %v", ourPath)
	log.Printf("Number of Transfers: %d", len(ourPath) + len(santaPath) - 2)
}
