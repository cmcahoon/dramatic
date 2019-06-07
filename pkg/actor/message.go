package actor

// CancelMessage will cause an actor to stop when processed.
type CancelMessage struct{}

// requestEnvelope lets the actor know that the sender would like a response
type requestEnvelope struct {
	Response *FutureTask
	Message  interface{}
}
