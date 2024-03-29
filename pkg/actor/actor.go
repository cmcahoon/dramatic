package actor

import (
	"sync"
)

type actor struct {
	name  string
	path  Path
	inbox <-chan interface{}
	group *sync.WaitGroup
	fn    ActorFn
}

func (a *actor) run() error {
	logger.Infow("starting actor", "actor_name", a.name, "actor_path", a.path.String())

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

			// Check if the sender wants a response
			var response *FutureTask
			message := msg

			envelope, ok := msg.(requestEnvelope)
			if ok {
				response = envelope.Response
				message = envelope.Message
			}

			// Pass the message and response future to the actor function
			err := a.fn(message, response)
			if err != nil {
				logger.Errorw("actor message handler failed", "error", err)
			}
		}
	}()

	return nil
}
