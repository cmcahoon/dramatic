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
	hd.MinLength = 8

	var err error
	hasher, err = hashids.NewWithData(hd)
	if err != nil {
		log.Fatal(err)
	}
}

// System is the root supervisor for multiple actors.
type System struct {
	name   string
	actors []*Ref
	group  sync.WaitGroup
}

// NewSystem creates a new System.
func NewSystem(name string) *System {
	return &System{
		name:   name,
		actors: make([]*Ref, 0),
		group:  sync.WaitGroup{},
	}
}

// NewActorFromFn will add, and immediately run, a new actor within the actor system.
func (s *System) NewActorFromFn(name string, fn ActorFn) *Ref {
	// Generate path
	actorCount = actorCount + 1
	id, err := hasher.Encode([]int{actorCount})
	if err != nil {
		logger.Fatal(err)
	}
	path := Path{id, name, ""}

	logger.Infow(
		"creating actor",
		"system_name", s.name,
		"actor_name", name,
		"actor_path", path.String(),
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
	ref := &Ref{inbox, path}
	s.actors = append(s.actors, ref)

	return ref
}

// NewActorFromStruct will add, and immediately run, a new actor within the actor system.
func (s *System) NewActorFromStruct(name string, actor Actor) *Ref {
	return s.NewActorFromFn(name, actor.Receive)
}

// Terminate will cancel all actors in the system and return. This call is blocking.
func (s *System) Terminate() error {
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
