package main

import (
	"github.com/cmcahoon/ligero/pkg/actor"
	"log"
	"time"
)

func main() {
	future := actor.NewFutureTask()

	// Simulate a long running task
	go func() {
		time.Sleep(5 * time.Second)
		err := future.Resolve("Done.")
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for the result
	log.Println("Waiting for future result...")

	result, err := future.GetResult()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}
