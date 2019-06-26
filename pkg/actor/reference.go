package actor

// Ref abstracts an actor away by providing an opaque reference to it. The reference can be used to publish
// messages to the actor.
type Ref struct {
	inbox chan<- interface{}
	path  Path
}

// Publish will publish any arbitrary message to the actor's inbox.
func (r *Ref) Publish(msg interface{}) {
	r.inbox <- msg
}

// Request kicks off a request-response cycle. Unlike `Publish`, this function expects a response from the target actor.
func (r *Ref) Request(msg interface{}) Future {
	future := NewFutureTask()

	// Wrap the message in an envelop with the sender information
	envelope := requestEnvelope{
		Response: future,
		Message:  msg,
	}
	r.inbox <- envelope
	return future
}
