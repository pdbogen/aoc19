package main

import (
	"github.com/pdbogen/aoc19/intcode"
	"log"
)

type Packet struct{ X, Y int }

type Node struct {
	Input      chan<- int
	Output     <-chan int
	InputReady <-chan struct{}
	Program    []int
	Queue      []Packet
	NeedInput  bool
}

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d character program %v", len(prog), prog)

	var nodes []*Node
	var NAT *Packet
	var lastNatY = -1
	for i := 0; i < 50; i++ {
		in := make(chan int)
		out := make(chan int)
		inputReady := make(chan struct{})
		node := &Node{
			Input:      in,
			Output:     out,
			InputReady: inputReady,
			Program:    make([]int, len(prog)),
		}
		copy(node.Program, prog)
		go func(i int, node *Node) {
			log.Printf("%d: starting", i)
			_, err := intcode.ExecuteInteractive(node.Program, in, inputReady, out)
			log.Printf("%d: exited: %v", i, err)
		}(i, node)
		nodes = append(nodes, node)
	}

	for i, node := range nodes {
		<-node.InputReady
		node.Input <- i
	}

	for {
		idle := 0
		for _, node := range nodes {
			if node.NeedInput {
				idle++
				continue
			}
			select {
			case <-node.InputReady:
				node.NeedInput = true
				idle++
			default:
			}
		}
		log.Printf("%d nodes waiting for input", idle)
		for i, node := range nodes {
			select {
			case dest, ok := <-node.Output:
				if !ok {
					log.Fatalf("%d: output channel closed", i)
				}
				pkt := &Packet{X: <-node.Output, Y: <-node.Output}
				log.Printf("%d: send %v to %d", i, pkt, dest)
				if dest == 255 {
					NAT = pkt
				} else {
					nodes[dest].Queue = append(nodes[dest].Queue, *pkt)
				}
			default:
			}
		}
		for _, node := range nodes {
			if !node.NeedInput {
				continue
			}
			if len(node.Queue) == 0 {
				continue
			}
			idle--
			node.NeedInput = false
			node.Input <- node.Queue[0].X
			<-node.InputReady
			node.Input <- node.Queue[0].Y
			node.Queue = node.Queue[1:]
		}
		if idle == 50 {
			if NAT == nil {
				for _, node := range nodes {
					node.Input <- -1
					node.NeedInput = false
				}
			} else {
				nodes[0].Input <- NAT.X
				<-nodes[0].InputReady
				nodes[0].Input <- NAT.Y
				nodes[0].NeedInput = false
				if NAT.Y == lastNatY {
					log.Fatalf("dupe Y %d", NAT.Y)
				}
				lastNatY = NAT.Y
				log.Printf("255: sending NAT packet %v to 0", NAT)
			}
		}
		//time.Sleep(10*time.Millisecond)
	}
}
