package ligero

// ActorRef abstracts an actor away by providing an opaque reference to it. The reference can be used to publish
// messages to the actor.
type ActorRef struct {
	inbox chan<- interface{}
}

// Publish will publish any arbitrary message to the actor's inbox.
func (a *ActorRef) Publish(msg interface{}) {
	a.inbox <- msg
}
