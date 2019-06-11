package main

import (
	"errors"
	"fmt"
	"github.com/cmcahoon/ligero/pkg/actor"
	"log"
)

type depositMessage struct {
	amount uint64
}

type withdrawlMessage struct {
	amount uint64
}

type getBalanceMessage struct{}

type accountActor struct {
	balance uint64
}

func (a *accountActor) Receive(msg interface{}, response *actor.FutureTask) error {
	switch typedMsg := msg.(type) {
	case depositMessage:
		a.balance = a.balance + typedMsg.amount
	case withdrawlMessage:
		if a.balance < typedMsg.amount {
			fmt.Println("overdraft protection kicked in - withdrawl blocked")
			return nil
		}
		a.balance = a.balance - typedMsg.amount
	case getBalanceMessage:
		if response == nil {
			return errors.New("response was nil")
		}
		err := response.Resolve(a.balance)
		if err != nil {
			return errors.New(err.Error())
		}
	default:
		return errors.New("unsupported message")
	}

	return nil
}

func main() {
	system := actor.NewActorSystem("system")

	account := accountActor{0}
	accountRef := system.NewActorFromStruct("account", &account)

	accountRef.Publish(depositMessage{100})   // 100
	accountRef.Publish(withdrawlMessage{25})  // 75
	accountRef.Publish(depositMessage{100})   // 175
	accountRef.Publish(withdrawlMessage{200}) // Overdraft!

	balance, err := accountRef.Request(getBalanceMessage{}).GetResult()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Balance is: %d", balance.(uint64))

	err = system.Terminate()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
