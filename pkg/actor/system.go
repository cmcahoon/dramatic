package actor

import (
	"github.com/speps/go-hashids"
	"log"
	"sync"
)

var actorCount = 0
var hasher *hashids.HashID

func init() {
	hd := hashids.NewData()
	hd.Salt = "6B7E6B29-049D-48A6-A414-2B8D1A821EE3"

	var err error
	hasher, err = hashids.NewWithData(hd)
	if err != nil {
		log.Fatal(err)
	}
}

// ActorSystem is the root supervisor for multiple actors.
type ActorSystem struct {
	name   string
	actors []*ActorRef
	group  sync.WaitGroup
}

// NewActorSystem creates a new ActorSystem.
func NewActorSystem(name string) *ActorSystem {
	return &ActorSystem{
		name:   name,
		actors: make([]*ActorRef, 0),
		group:  sync.WaitGroup{},
	}
}

// NewActorFromFn will add, and immediately run, a new actor within the actor system.
func (s *ActorSystem) NewActorFromFn(name string, fn ActorFn) *ActorRef {
	// Generate path
	actorCount = actorCount + 1
	id, err := hasher.Encode([]int{actorCount})
	if err != nil {
		logger.Fatal(err)
	}
	path := "/" + name + "/" + id

	logger.Infow(
		"creating actor",
		"system_name", s.name,
		"actor_name", name,
		"actor_path", path,
	)

	// Create the inbox and the actor
	inbox := make(chan interface{})
	actor := &actor{
		name:  name,
		path:  path,
		inbox: inbox,
		group: &s.group,
		fn:    fn,
	}

	// Run the actor
	err = actor.run()
	if err != nil {
		log.Fatal(err)
	}

	// Create and store the reference
	ref := &ActorRef{inbox, path}
	s.actors = append(s.actors, ref)

	return ref
}

// NewActorFromStruct will add, and immediately run, a new actor within the actor system.
func (s *ActorSystem) NewActorFromStruct(name string, actor Actor) *ActorRef {
	return s.NewActorFromFn(name, actor.Receive)
}

// Terminate will cancel all actors in the system and return. This call is blocking.
func (s *ActorSystem) Terminate() error {
	logger.Infow(
		"terminating actor system",
		"system_name", s.name,
	)

	for _, ref := range s.actors {
		ref.Publish(CancelMessage{})
	}

	s.group.Wait()
	return nil
}
