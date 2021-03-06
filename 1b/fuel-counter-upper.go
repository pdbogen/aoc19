package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	moduleMasses, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("could not read input: %+v", err)
	}
	massList := strings.Split(strings.TrimSpace(string(moduleMasses)), "\n")
	log.Printf("got %d modules", len(massList))
	totalFuel := 0
	for _, massStr := range massList {
		mass, err := strconv.Atoi(massStr)
		if err != nil {
			log.Fatalf("could not convert mass %q to integer: %+v", massStr, err)
		}
		addlFuel := mass
		for {
			addlFuel = addlFuel/3 - 2
			if addlFuel <= 0 {
				break
			}
			totalFuel += addlFuel
		}
	}
	log.Printf("Required fuel: %d", totalFuel)
}
