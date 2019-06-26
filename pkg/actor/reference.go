package actor

// ActorRef abstracts an actor away by providing an opaque reference to it. The reference can be used to publish
// messages to the actor.
type ActorRef struct {
	inbox chan<- interface{}
	path  Path
}

// Publish will publish any arbitrary message to the actor's inbox.
func (a *ActorRef) Publish(msg interface{}) {
	a.inbox <- msg
}

// Request kicks off a request-response cycle. Unlike `Publish`, this function expects a response from the target actor.
func (a *ActorRef) Request(msg interface{}) Future {
	future := NewFutureTask()

	// Wrap the message in an envelop with the sender information
	envelope := requestEnvelope{
		Response: future,
		Message:  msg,
	}
	a.inbox <- envelope
	return future
}
