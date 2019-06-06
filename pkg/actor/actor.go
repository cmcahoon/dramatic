package actor

import (
	"log"
	"sync"
)

type actor struct {
	name  string
	inbox <-chan interface{}
	group *sync.WaitGroup
	fn    ActorFn
}

func (a *actor) run() error {
	logger.Infow("starting actor", "actor_name", a.name)

	go func() {
		a.group.Add(1)
		defer a.group.Done()

		for {
			msg := <-a.inbox

			// Handle a cancel before passing on to the actor function. This is transparent to the implementor.
			_, ok := msg.(CancelMessage)
			if ok {
				logger.Infow("cancel received, stopping actor", "actor_name", a.name)
				return
			}

			// Otherwise, pass the message to the actor function
			err := a.fn(msg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	return nil
}
