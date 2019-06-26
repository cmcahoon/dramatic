package actor

// Path is a unique string that identifies an actor and specifies it's ancestry.
//
// For example:  /foo/pz4h1
//   	The `ancestry` is '/'. This means the actor has no ancestors -- it's a root actor.
//		The `name` is 'foo'. There can be multiple actors with the same name.
//		The `id` is 'pz4h1' which uniquely identifies the actor.
type Path struct {
	id       string
	name     string
	ancestry string
}

func (p *Path) String() string {
	return p.ancestry + "/" + p.name + "/" + p.id
}
