package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	producerCount        = 4
	consumerCount        = 2
	bufferSize           = 10
	simulationCycleCount = 10
)

type Widget struct {
}

func produce(channel chan Widget, producerId int) {
	for {
		widgetsCount := rand.Intn(3) + 1
		fmt.Printf("Producing widets %d by producer %d\n", widgetsCount, producerId)
		for widgetsCount != 0 {
			channel <- Widget{}
			widgetsCount--
		}
	}
}

func consume(channels []chan Widget, waitGroup *sync.WaitGroup, consumerId int) {
	currentProducer := 0
	widgets := make([]Widget, 0, bufferSize)
	simulationCount := 0
	for {
		widgetsCount := rand.Intn(3) + 1
		widgetsConsumed := 0

		for widgetsConsumed < widgetsCount && len(widgets) < cap(widgets) {
			widget := <-channels[currentProducer]
			widgets = append(widgets, widget)
			widgetsConsumed++
		}
		fmt.Printf("Consumed widets %d, from producer %d\n", widgetsConsumed, currentProducer+1)
		currentProducer = (currentProducer + 1) % producerCount

		if len(widgets) == cap(widgets) {
			simulationCount++
			fmt.Printf("discarding %+v\n", consumerId)
			widgets = make([]Widget, 0, bufferSize)
			timeToSleep := rand.Intn(5) + 1
			time.Sleep(time.Duration(timeToSleep) * time.Second)
			if simulationCount == simulationCycleCount {
				break
			}
		}
	}
	waitGroup.Done()
}

func main() {
	var channels = make([]chan Widget, producerCount)

	waitGroup := sync.WaitGroup{}
	for i := 0; i < producerCount; i++ {
		channels[i] = make(chan Widget, bufferSize)
		go produce(channels[i], i+1)
	}

	for i := 0; i < consumerCount; i++ {
		waitGroup.Add(1)
		go consume(channels, &waitGroup, i+1)
	}
	waitGroup.Wait()
}
