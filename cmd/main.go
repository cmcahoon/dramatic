package main

import (
	"errors"
	"github.com/cmcahoon/ligero/pkg/actor"
	"log"
)

func main() {
	system := actor.NewActorSystem("system")
	doubleActor := system.NewActor("double", func(msg interface{}) error {
		switch typedMsg := msg.(type) {
		case int:
			log.Printf("Double of %d is %d", typedMsg, typedMsg*2)
		default:
			return errors.New("unsupported message")
		}

		return nil
	})

	doubleActor.Publish(2)

	err := system.Terminate()
	if err != nil {
		log.Fatal(err)
	}
}
