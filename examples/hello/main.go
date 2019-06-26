package main

import (
	"errors"
	"fmt"
	"github.com/cmcahoon/ligero/pkg/actor"
)

func main() {
	system := actor.NewSystem("system")
	helloActor := system.NewActorFromFn("hello", func(msg interface{}, _ *actor.FutureTask) error {
		switch typedMsg := msg.(type) {
		case string:
			fmt.Printf("Hello, %s!\n", typedMsg)
		default:
			return errors.New("unsupported message")
		}

		return nil
	})

	helloActor.Publish("Actor")

	err := system.Terminate()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
