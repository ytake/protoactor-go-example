package task

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/event"
	"math/rand"
)

const (
	FailBeforeProcessing = iota
	FailAfterProcessing
	ProcessSuccessfully
)

type (
	Behavior int
	Account struct {
		Name               string
		ServiceUptime      float64
		RefusalProbability float64
		BusyProbability    float64
		ProcessedMessages  map[*actor.PID]interface{}
		Balance            int
	}
)

// Receive 銀行アカウント送金
func (a *Account) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *event.Credit :
		a.AdjustBalance(context, msg.ReplyTo, msg.Amount)
	}
}

func (a *Account) Busy() bool {
	comparison := rand.Float64() * 100
	return comparison <= a.BusyProbability
}

func (a *Account) RefusePermanently() bool {
	comparison := rand.Float64() * 100
	return comparison <= a.RefusalProbability
}

// Failure 失敗
func (a *Account) Failure(context actor.Context, replyTo *actor.PID) {
	context.Send(replyTo, event.InternalServerError{})
}

func (a *Account) DetermineProcessingBehavior() Behavior {
	comparison := rand.Float64() * 100
	if comparison > a.ServiceUptime {
		if rand.Float64()*100 > 50 {
			return FailBeforeProcessing
		}
		return FailAfterProcessing
	}
	return ProcessSuccessfully
}

func (a *Account) AlreadyProcessed(replyTo *actor.PID) bool {
	_, ok := a.ProcessedMessages[replyTo]
	return ok
}

func (a *Account) AdjustBalance(context actor.Context, replyTo *actor.PID, amount int) {
	if a.RefusePermanently() {
		a.ProcessedMessages[replyTo] = &event.Refused{}
		context.Send(replyTo, &event.Refused{})
		return
	}
	if a.Busy() {
		context.Send(replyTo, &event.ServiceUnavailable{})
	}
}
