package actor

// Actor is any object than can receive an arbitrary message.
type Actor interface {
	Receive(msg interface{}, response *FutureTask) error
}

// ActorFn is a function that will be called on each message received in the inbox.
type ActorFn func(msg interface{}, response *FutureTask) error
