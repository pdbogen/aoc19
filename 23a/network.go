package main

import (
	"github.com/pdbogen/aoc19/intcode"
	"log"
	"sync/atomic"
	"time"
)

var idle int32

type Packet struct {
	X, Y int
}

func QueueManager(input chan<- int, queue <-chan Packet) {
	log.Print("Queue manager starting")
	q := []Packet{{-1, -1}}
	for {
		select {
		case i := <-queue:
			if q[0].X == -1 {
				atomic.AddInt32(&idle, 1)
				q[0] = i
			} else {
				q = append(q, i)
			}
		case input <- q[0].X:
			if q[0].X == -1 {
				continue
			}
			input <- q[0].Y
			if len(q) == 1 {
				atomic.AddInt32(&idle, 0)
				q[0].X = -1
			} else {
				q = q[1:]
			}
		}
	}
}

func PacketReceiver(output <-chan int, queues map[int]chan Packet) {
	log.Print("Packet receiver starting")
	for address := range output {
		if address == 255 {
			log.Print("packet for 255: ", Packet{<-output, <-output})
			continue
		}

		pkt := Packet{<-output, <-output}
		queues[address] <- pkt
		log.Printf("queued %v for %d (idle=%d)", pkt, address, atomic.LoadInt32(&idle))
	}
}

func NAT(in <-chan Packet, out chan<- Packet) {
	p := Packet{-1, 0}
	for {
		select {
		case p = <-in:
		default:
			if atomic.LoadInt32(&idle) == 50 && p.X != -1 {
				out <- p
				p.X = -1
			}
		}
	}
}

func main() {
	prog, err := intcode.LoadFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d character program", len(prog))

	inputs := map[int]chan int{}
	outputs := map[int]chan int{}
	queues := map[int]chan Packet{}

	for i := 0; i < 50; i++ {
		log.Printf("starting node %d", i)
		inputs[i] = make(chan int)
		outputs[i] = make(chan int)
		queues[i] = make(chan Packet)
		go QueueManager(inputs[i], queues[i])
		go intcode.Execute(prog, inputs[i], outputs[i])
		inputs[i] <- i
	}
	for i := 0; i < 50; i++ {
		go PacketReceiver(outputs[i], queues)
	}
	queues[255] = make(chan Packet)
	go NAT(queues[255], queues[0])

	for {
		time.Sleep(time.Minute)
	}
}
