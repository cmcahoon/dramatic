# dramatic - Actors for Go.
[![CircleCI](https://circleci.com/gh/cmcahoon/dramatic.svg?style=svg)](https://circleci.com/gh/cmcahoon/dramatic)
&nbsp;
[![codecov](https://codecov.io/gh/cmcahoon/dramatic/branch/master/graph/badge.svg)](https://codecov.io/gh/cmcahoon/dramatic)



## First Things First

**This library is not production-quality.**

I wrote this project to further explore the actor model. I've used [Akka](https://akka.io/) with Java -- highly recommended if you are looking for a full fledged Actor deployment -- and experimented with [Kotlin Actors](https://kotlinlang.org/docs/reference/coroutines/shared-mutable-state-and-concurrency.html#actors). Many ideas that made their way into `dramatic` were stolen from these excellent projects. To take this to the next level, I will occasionally link to Akka documentation because they have spent a lot of effort explaining the actor model and why you should consider them.

## Why Actors?

Remember when life was easy, the world was peaceful, and your code ran in a single thread? Well, not to be downer, but those days are over. Software is distributed now, processors can run many threads in parallel, and our old way of programming hasn't kept up. We resort to complex locking strategies and manually debugging race conditions and deadlocks. Actors are not exactly new, but they a promising solution to these common problems.

The basic idea is that an actor owns state and that state can only be modified by that actor. Actors only modify state if they receive a message to do so. You can think of the messages much like using a message broker or work queues; they are FIFO. An actor, processing a single message at once, cannot have a race condition. The message queue between actors acts as the synchronization mechanism instead of locks, mutexes, semaphores, etc.

You'll see this idea in Go, for example, where they famously state "Do not communicate by sharing memory; share memory by communicating." They encourage the use of channels -- essentially a queue -- as a sychronization mechanism. They also realized that the current model of handling concurrency is inadequate.

Obviously there is more going on here, so I'll leave you with some great [documentation on actors](https://doc.akka.io/docs/akka/current/guide/actors-motivation.html) from the real pros at Akka.

## Use

### Import Location

The project can be imported from `github.com/cmcahoon/dramatic/pkg/actor`.

### Create an Actor System

All actors belong to an actor system. In `dramatic`, the actor system is mostly just a container for actors and provides the ability to terminate them all at once.

```go
system := actor.NewSystem("system")
```

### Create an Actor

There are two ways to implement and actor, by function or by implementing the `Actor` interface. Regardless of the implementation method, the code looks almost similar:

```go
// This creates `helloActor` from a function that takes a message of arbitrary type.
helloActorRef := system.NewActorFromFn("hello", func(msg interface{}, _ *actor.FutureTask) error {
    switch typedMsg := msg.(type) {
    case string:
        fmt.Printf("Hello, %s!\n", typedMsg)
    default:
        return errors.New("unsupported message")
    }

    return nil
})

// This creates `helloActor` by implementing the `Actor` interface.
type HelloActor struct {}

func (a *HelloActor) Receive(msg interface{}, _ *actor.FutureTask) error {
        switch typedMsg := msg.(type) {
    case string:
        fmt.Printf("Hello, %s!\n", typedMsg)
    default:
        return errors.New("unsupported message")
    }

    return nil
}

helloActorRef := system.NewActorFromStruct("hello", &HelloActor{})
```

> What's with that `actor.FutureTask` parameter? We'll get there in a minute.

### Send Messages

Now that we have our actor running, how do we send messages to it? Each time you add an actor to the actor system, you receive an actor reference back. This allows you to send messages to that specific actor.

```go
helloActorRef.Publish("Actor")
```

This will send a message of type `string` with value "Actor" to the "hello" actor. The `publish` function will return once the message has been placed in the channel; it **will not** wait for a response.

### Requesting Responses from Actors

Sometimes you want to wait for a response from an actor. This is where the `actor.FutureTask` parameter comes in. Let's revisit creating an actor, and this time name the second parameter:

```go
type getBalanceMessage struct{}

accountActorRef := system.NewActorFromFn("account", func(msg interface{}, result *actor.FutureTask) error {
	switch typedMsg := msg.(type) {
	case getBalanceMessage:
		err := response.Resolve(1_000_000)
		if err != nil {
			return errors.New(err.Error())
		}
	default:
		return errors.New("unsupported message")
	}

    return nil
})

// Instead of using `publish`, use `request`. This will return a future you can wait on. Be aware, `GetResult` on the
// future will block the calling thread.
balance, err := accountActorRef.Request(getBalanceMessage{}).GetResult()
if err != nil {
    log.Fatal(err)
}
log.Printf("Balance is: %d", balance.(uint64))
```
